package network

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const wsBufferSize = 1024

func wsOriginChecker(_ *http.Request) bool {
	return true
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  wsBufferSize,
	WriteBufferSize: wsBufferSize,
	CheckOrigin:     wsOriginChecker,
}

func WsUpgrade(w http.ResponseWriter, r *http.Request, gameId, playerId uuid.UUID, cm ConnectionManager) (*WsConnection, error) {
	c, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %s", err)
		return nil, err
	}

	conn := cm.NewConnection(c, gameId, playerId)
	return conn, nil
}

type GameMessage struct {
	Message  string
	GameId   uuid.UUID
	PlayerId uuid.UUID
}

type WsConnection struct {
	conn         NetworkConnection
	GameId       uuid.UUID
	PlayerId     uuid.UUID
	JoinTime     time.Time
	LastPingTime time.Time
	lock         sync.Mutex
}

func (wsconn *WsConnection) Send(msg []byte) error {
	wsconn.lock.Lock()
	defer wsconn.lock.Unlock()
	return wsconn.conn.Send(msg)
}

func (wsconn *WsConnection) Close() {
	wsconn.lock.Lock()
	defer wsconn.lock.Unlock()
	wsconn.conn.Close()
}

func (gcm *IntegratedConnectionManager) NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection {
	c := &WsConnection{conn: &WebsocketConnection{Conn: conn},
		PlayerId:     playerId,
		GameId:       gameId,
		JoinTime:     time.Now(),
		LastPingTime: time.Now(),
	}
	gcm.RegisterConnection(gameId, playerId, c)
	return c
}

// To be called after registering the connection, this will listen to the
// websocket traffic on a loop and handle it
func (c *WsConnection) listenAndHandle() error {
	gid := c.GameId
	// pid := c.PlayerId
	game, err := GameRepo.GetGame(gid)
	if err != nil {
		return err
	}

	state := RpcOnJoinMsg(RpcOnJoinMsg{State: game.StateInfo()})
	initialState, err := EncodeRpcMessage(state)
	if err != nil {
		return err
	}

	err = c.Send(initialState)
	if err != nil {
		return err
	}

	// Start listening and handling
	for {
		msg, err := c.conn.Receive()
		if err != nil {
			c.Close()
			return errors.New("Cannot read from websocket")
		}

		log.Printf("Got a message: %s", string(msg))
		var message RpcMessage
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Printf("Error unmarshalling message: %s; closing connection for gid %s pid %s", err, c.GameId, c.PlayerId)
			c.Close()
			continue
		}
	}
}

func (c *WsConnection) ListenAndHandle(g *IntegratedConnectionManager) {
	err := c.listenAndHandle()
	if err != nil {
		log.Printf("Error whilst handling websocket connection %s", err)
		c.Close()
		g.UnregisterConnection(c.GameId, c.PlayerId)
	}
}
