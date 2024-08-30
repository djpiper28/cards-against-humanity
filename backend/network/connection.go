package network

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
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

func WsUpgrade(w http.ResponseWriter, r *http.Request, gameId, playerId uuid.UUID, cm *ConnectionManager) (*WsConnection, error) {
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
	// Used to know if a ping has been sent
	PingFlag bool
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
	wsconn.Connected = false
}

const pingInterval = 5 * time.Second
const pingTimeout = 2 * pingInterval

func (gcm *ConnectionManager) NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection {
	c := &WsConnection{conn: &WebsocketConnection{Conn: conn},
		PlayerId:     playerId,
		GameId:       gameId,
		JoinTime:     time.Now(),
		PingFlag:     false,
		Connected:    true,
		LastPingTime: time.Now(),
	}
	gcm.RegisterConnection(gameId, playerId, c)

	// Ping timeout handler
	go func() {
		time.Sleep(pingInterval)
		for {
			err := func() error {
				c.lock.Lock()
				defer c.lock.Unlock()

				if !c.Connected {
					return errors.New("Connection is not connected")
				}

				if time.Since(c.LastPingTime) > pingInterval && !c.PingFlag {
					pingMessage, err := EncodeRpcMessage(RpcPingMsg{})
					if err != nil {
						return err
					}
					c.PingFlag = true

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
const PingCommand = "Ping"

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

	// TODO: do not block on this call - this will cause bugs
	GlobalConnectionManager.Broadcast(gid, message)

	// Send the initial state
	state := RpcOnJoinMsg(RpcOnJoinMsg{State: game.StateInfo(c.PlayerId)})
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
				handler = PingCommand
				c.lock.Lock()
				defer c.lock.Unlock()

				c.LastPingTime = time.Now()
				c.PingFlag = false
				return nil
			},
			StartGameHandler: func() error {
				handler = "Start Game"
				game, err := gameRepo.Repo.GetGame(gid)
				if err != nil {
					return errors.New("Cannot find the game")
				}

				if game.GameOwnerId != c.PlayerId {
					return errors.New("Only the game owner can start the game")
				}

				info, err := gameRepo.Repo.StartGame(gid)
				if err != nil {
					return err
				}

				totalPlays := 0
				for _, plays := range info.PlayersPlays {
					if len(plays) > 0 {
						totalPlays++
					}
				}

				for playerId, hand := range info.PlayerHands {
					handCopy := make([]gameLogic.WhiteCard, len(hand))
					for i, val := range hand {
						handCopy[i] = *val
					}

					playsCopy := make([]gameLogic.WhiteCard, 0)
					for i, val := range info.PlayersPlays[playerId] {
						playsCopy[i] = *val
					}

					roundInfo := RpcRoundInformationMsg{
						CurrentCardCzarId: info.CurrentCardCzarId,
						YourHand:          handCopy,
						RoundNumber:       info.RoundNumber,
						BlackCard:         *info.CurrentBlackCard,
						YourPlays:         playsCopy,
						TotalPlays:        totalPlays,
					}

					encodedMessage, err := EncodeRpcMessage(roundInfo)
					if err != nil {
						return err
					}

					go GlobalConnectionManager.SendToPlayer(c.GameId, playerId, encodedMessage)
				}
				return nil
			},
			PlayCardsHandler: func(msg RpcPlayCardsMsg) error {
				handler = "Play Cards"

				info, err := gameRepo.Repo.PlayerPlayCards(c.GameId, c.PlayerId, msg.CardIds)
				if err != nil {
					return err
				}

				cardPlayedMsg := RpcOnCardPlayedMsg{
					PlayerId: c.PlayerId,
				}

				broadcastMessage, err := EncodeRpcMessage(cardPlayedMsg)
				if err != nil {
					return err
				}

				go GlobalConnectionManager.Broadcast(gid, broadcastMessage)
				if info.MovedToNextCardCzarPhase {
					go GlobalConnectionManager.MoveToCzarJudgingPhase(gid, info.CzarJudingPhaseInfo)
				}
				return nil
			},
			CzarSelectCardHandler: func(msg RpcCzarSelectCardMsg) error {
				handler = "Czar Selects A Card"

				res, err := gameRepo.Repo.CzarSelectsCard(c.GameId, c.PlayerId, msg.Cards)
				if err != nil {
					return err
				}

				if res.GameEnded {
					// TODO: end the game lol
				}

				var wg sync.WaitGroup
				for pid, hand := range res.Hands {
					wg.Add(1)
					go func(pid uuid.UUID, hand []*gameLogic.WhiteCard) {
						defer wg.Done()
						msg := RpcOnWhiteCardPlayPhase{YourHand: hand,
							BlackCard:  res.NewBlackCard,
							CardCzarId: res.NewCzarId}
						encodedMsg, err := EncodeRpcMessage(msg)
						if err != nil {
							logger.Logger.Error("Cannot encode message to send to player")
						}

						go GlobalConnectionManager.SendToPlayer(c.GameId, pid, encodedMsg)
					}(pid, hand)
				}

				wg.Wait()
				return nil 
			},
		})

		microSeconds := time.Since(startTime).Microseconds()
		go gameRepo.AddCommandExecuted(int(time.Since(startTime).Microseconds()))

		if handler != PingCommand {
			logger.Logger.Infof("Command Handler \"%s\" | %s | %dµs",
				handler,
				gid,
				microSeconds)
		}

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

func (c *WsConnection) ListenAndHandle(g *ConnectionManager) {
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
