package network

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
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
		logger.Logger.Error("Failed to set websocket upgrade",
			"err", err)
		return nil, err
	}

	go gameRepo.AddWsConnection()

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
	// Used to terminate the ping thread
	Connected bool
	lock      sync.Mutex
}

func (wsconn *WsConnection) Send(msg []byte) error {
	wsconn.lock.Lock()
	defer wsconn.lock.Unlock()
	go gameRepo.AddMessageSent()
	return wsconn.conn.Send(msg)
}

func (wsconn *WsConnection) Close() {
	wsconn.lock.Lock()
	defer wsconn.lock.Unlock()
	wsconn.conn.Close()
}

const pingTimeout = 10 * time.Second
const pingInterval = 5 * time.Second

func (gcm *IntegratedConnectionManager) NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection {
	c := &WsConnection{conn: &WebsocketConnection{Conn: conn},
		PlayerId:     playerId,
		GameId:       gameId,
		JoinTime:     time.Now(),
		Connected:    true,
		LastPingTime: time.Now(),
	}
	gcm.RegisterConnection(gameId, playerId, c)

	// Ping timeout handler
	go func() {
		for {
			err := func() error {
				c.lock.Lock()
				defer c.lock.Unlock()

				if !c.Connected {
					return errors.New("Connection is not connected")
				}

				if time.Since(c.LastPingTime) > pingInterval {
					pingMessage, err := EncodeRpcMessage(RpcPingMsg{})
					if err != nil {
						return err
					}

					go c.Send(pingMessage)
					return nil
				} else if time.Since(c.LastPingTime) > pingTimeout {
					logger.Logger.Warn("Player timed out due to no ping response",
						"playerId", c.PlayerId,
						"gameId", c.GameId)
					go c.Close()
					return errors.New("Ping timeout")
				}
				return nil
			}()

			if err != nil {
				return
			}

			time.Sleep(time.Millisecond * 100)
		}
	}()
	return c
}

const UnknownCommand = "Unknown Command"

// To be called after registering the connection, this will listen to the
// websocket traffic on a loop and handle it
func (c *WsConnection) listenAndHandle() error {
	gid := c.GameId
	// pid := c.PlayerId
	game, err := gameRepo.Repo.GetGame(gid)
	if err != nil {
		return err
	}

	name, err := gameRepo.Repo.GetPlayerName(gid, c.PlayerId)
	if err != nil {
		logger.Logger.Error("Cannot get the player's name", "err", err)
		name = "Error"
	}

	onPlayerJoinmsg := RpcOnPlayerJoinMsg{
		Id:   c.PlayerId,
		Name: name,
	}

	message, err := EncodeRpcMessage(onPlayerJoinmsg)
	if err != nil {
		logger.Logger.Error("Cannot encode the messages", "err", err)
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

		go gameRepo.AddCommandExecuted()

		handler := UnknownCommand
		startTime := time.Now()
		err = DecodeRpcMessage(msg, RpcCommandHandlers{
			ChangeSettingsHandler: func(msg RpcChangeSettingsMsg) error {
				handler = "Change Settings"
				game, err := gameRepo.Repo.GetGame(gid)
				if err != nil {
					return errors.New("Cannot find the game")
				}

				if game.GameOwnerId != c.PlayerId {
					return errors.New("You cannot change the settings as you are not the game owner")
				}

				err = gameRepo.Repo.ChangeSettings(gid, msg.Settings)
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
			PingHandler: func() error {
				c.lock.Lock()
				defer c.lock.Unlock()

				c.LastPingTime = time.Now()
				return nil
			},
		})

		endTime := time.Now()
		microSeconds := endTime.Sub(startTime).Microseconds()

		logger.Logger.Infof("Command Handler \"%s\" | %s | %dµs",
			handler,
			gid,
			microSeconds)

		if handler == UnknownCommand {
			go gameRepo.AddUnknownCommand()
		}

		if err != nil {
			logger.Logger.Error("Error processing message",
				"err", err,
				"gameId", c.GameId,
				"playerId", c.PlayerId)
			go gameRepo.AddCommandFailed()

			var message RpcCommandErrorMsg
			message.Reason = err.Error()
			encodedMessage, err := EncodeRpcMessage(message)
			if err != nil {
				logger.Logger.Error("Cannot encode the error message",
					"err", err)
				continue
			}

			go c.Send(encodedMessage)
		}
	}
}

func (c *WsConnection) ListenAndHandle(g *IntegratedConnectionManager) {
	err := c.listenAndHandle()
	if err != nil {
		logger.Logger.Error("Error whilst handling websocket connection",
			"err", err)
		go gameRepo.AddWsError()
		go gameRepo.RemoveWsConnection()
		c.Close()
		g.UnregisterConnection(c.GameId, c.PlayerId)
	}
}
