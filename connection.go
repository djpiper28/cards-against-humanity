package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const wsBufferSize = 1024

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  wsBufferSize,
	WriteBufferSize: wsBufferSize,
}

func WsUpgrade(w http.ResponseWriter, r *http.Request, playerId, gameId uuid.UUID) (*WsConnection, error) {
	c, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %s", err)
		return nil, err
	}

	conn := NewConnection(c, playerId, gameId)
	go conn.Process()
	return conn, nil
}

type GameMessage struct {
	Message  string
	GameId   uuid.UUID
	PlayerId uuid.UUID
}

type WsConnection struct {
	Conn         *websocket.Conn
	PlayerId     uuid.UUID
	GameID       uuid.UUID
	JoinTime     time.Time
	LastPingTime time.Time
	WsRecieve    chan GameMessage
	WsBroadcast  chan string
	shutdown     chan bool
}

func NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection {
	return &WsConnection{Conn: conn,
		PlayerId:     playerId,
		GameID:       gameId,
		JoinTime:     time.Now(),
		LastPingTime: time.Now(),
		WsRecieve:    make(chan GameMessage),
		WsBroadcast:  make(chan string),
		shutdown:     make(chan bool),
	}
}

func (c *WsConnection) Process() {
	go func() {
		for {
			select {
			case <-c.shutdown:
				return
			case msg := <-c.WsBroadcast:
				c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
				globalConnectionManager.Close(c.GameID, c.PlayerId)
				return
			}
		}
	}()

	for {
		select {
		case <-c.shutdown:
			return
		default:
			msgType, msg, err := c.Conn.ReadMessage()
			if err != nil {
				log.Println(err)
				globalConnectionManager.Close(c.GameID, c.PlayerId)
				return
			}

			if msgType == websocket.TextMessage {
				c.WsRecieve <- GameMessage{Message: string(msg), GameId: c.GameID, PlayerId: c.PlayerId}
			}
		}
	}
}

func (c *WsConnection) Close() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Printf("Closing the websocket caused an issue %s", err)
			}
		}()

		c.Conn.Close()
		c.shutdown <- true

		close(c.WsBroadcast)
		close(c.WsRecieve)
	}()
}

type GameConnection struct {
	// Maps a player ID to a ws connection
	PlayerConnectionMap map[uuid.UUID]*WsConnection
	Broadcast           chan string
	lock                sync.Mutex
}

func (gc *GameConnection) Close(playerId uuid.UUID) error {
	gc.lock.Lock()
	defer gc.lock.Unlock()

	conn, found := gc.PlayerConnectionMap[playerId]
	if !found {
		return errors.New("Cannot remove a wensocket that has been closed already")
	}

	conn.Close()
	delete(gc.PlayerConnectionMap, playerId)
	return nil
}

var globalConnectionManager GlobalConnectionManager

type GlobalConnectionManager struct {
	// Maps a game ID to the game connection pool
	GameConnectionMap map[uuid.UUID]*GameConnection
	lock              sync.Mutex
}

func InitGlobalConnectionManager() {
	globalConnectionManager = GlobalConnectionManager{GameConnectionMap: make(map[uuid.UUID]*GameConnection)}
}

func (g *GlobalConnectionManager) Close(gameId, playerId uuid.UUID) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	gameConnectionMgr, found := g.GameConnectionMap[gameId]
	if !found {
		return errors.New("Cannot close a player connection from a game that has been closed")
	}

	return gameConnectionMgr.Close(playerId)
}
