package network

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

	name, err := GameRepo.GetPlayerName(gid, c.PlayerId)
	if err != nil {
		log.Printf("Cannot get the player's name: %s", err)
		name = "Error"
	}

	onPlayerJoinmsg := RpcOnPlayerJoinMsg{
		Id:   c.PlayerId,
		Name: name,
	}

	message, err := EncodeRpcMessage(onPlayerJoinmsg)
	if err != nil {
		log.Printf("Cannot encode the message: %s", err)
	}

	GlobalConnectionManager.Broadcast(gid, message)

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

		handler := "Unknown"
		startTime := time.Now()
		err = DecodeRpcMessage(msg, RpcCommandHandlers{
			ChangeSettingsHandler: func(msg RpcChangeSettingsMsg) error {
				handler = "Change Settings"
				game, err := GameRepo.GetGame(gid)
				if err != nil {
					return errors.New("Cannot find the game")
				}

				if game.GameOwnerId != c.PlayerId {
					return errors.New("You cannot change the settings as you are not the game owner")
				}

				err = GameRepo.ChangeSettings(gid, msg.Settings)
				if err != nil {
					return err
				}

				broadcastMessage, err := EncodeRpcMessage(msg)
				if err != nil {
					return err
				}

				go GlobalConnectionManager.Broadcast(gid, broadcastMessage)
				return nil
			},
		})

		endTime := time.Now()
		microSeconds := endTime.Sub(startTime).Microseconds()

		log.Printf("Command Handler \"%s\" | %s | %dµs", handler, gid, microSeconds)

		if err != nil {
			log.Printf("Error processing message: %s; for gid %s pid %s", err, c.GameId, c.PlayerId)

			var message RpcCommandErrorMsg
			message.Reason = err.Error()
			encodedMessage, err := EncodeRpcMessage(message)
			if err != nil {
				log.Printf("Cannot encode the error message: %s", err)
				continue
			}

			go c.Send(encodedMessage)
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
