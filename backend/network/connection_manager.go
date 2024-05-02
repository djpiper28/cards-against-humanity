package network

import (
	"errors"
	"log"
	"sync"

	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameConnection struct {
	// Maps a player to the connection
	playerConnectionMap map[uuid.UUID]*WsConnection
}

// Manages all of the connections
var GlobalConnectionManager = &IntegratedConnectionManager{GameConnectionMap: make(map[uuid.UUID]*GameConnection)}

type IntegratedConnectionManager struct {
	// Maps a game ID to the game connection pool
	GameConnectionMap map[uuid.UUID]*GameConnection
	lock              sync.Mutex
}

type ConnectionManager interface {
	RegisterConnection(gameId, playerId uuid.UUID, connection *WsConnection)
	NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection
	Close(gameId, playerId uuid.UUID) error
	RemovePlayer(gameId, playerId uuid.UUID) error
	RemoveGame(gameId uuid.UUID) error
}

func (g *IntegratedConnectionManager) Close(gameId, playerId uuid.UUID) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	game, found := g.GameConnectionMap[gameId]
	if !found {
		return errors.New("Cannot close a player connection from a game that has been closed")
	}

	game.playerConnectionMap[playerId].Close()
	delete(game.playerConnectionMap, playerId)
	return nil
}

func (g *IntegratedConnectionManager) RegisterConnection(gameId, playerId uuid.UUID, connection *WsConnection) {
	g.lock.Lock()
	defer g.lock.Unlock()

	logger.Logger.Info("Registering player to game",
		"playerId", playerId,
		"gameId", gameId)
	game, found := g.GameConnectionMap[gameId]
	if !found {
		game = &GameConnection{make(map[uuid.UUID]*WsConnection)}
		g.GameConnectionMap[gameId] = game
		logger.Logger.Info("Registered game", "gameId", gameId)
	}

	playerConnection, foundPlayer := game.playerConnectionMap[playerId]
	if foundPlayer {
		// Do not trigger a disconnect event in the internal system, just hot-swap the connection
		playerConnection.conn.Close()
	}

	game.playerConnectionMap[playerId] = connection

	err := gameRepo.Repo.ConnectPlayer(gameId, playerId)
	if err != nil {
		log.Printf("Cannot tag player %s as connected to game %s", playerId, gameId)
	}

	go connection.ListenAndHandle(g)
}

func (g *IntegratedConnectionManager) UnregisterConnection(gameId, playerId uuid.UUID) {
	g.lock.Lock()
	defer g.lock.Unlock()

	logger.Logger.Info("Unregistering player to game",
		"playerId", playerId,
		"gameId", gameId)
	game, found := g.GameConnectionMap[gameId]
	if found {
		delete(game.playerConnectionMap, playerId)
	} else {
		logger.Logger.Error("Cannot unregister game %s as it cannot be found",
			"gameId", gameId)
	}

	err := gameRepo.Repo.DisconnectPlayer(gameId, playerId)
	if err != nil {
		logger.Logger.Error("Cannot tag player as disconnected from game",
			"playerId", playerId,
			"gameId", gameId)
	}

	onPlayerDisconnectMsg := RpcOnPlayerDisconnectMsg{
		Id: playerId,
	}

	message, err := EncodeRpcMessage(onPlayerDisconnectMsg)
	if err != nil {
		logger.Logger.Error("Cannot encode the message",
			"err", err)
		return
	}

	go g.Broadcast(gameId, message)
}

// Blocking call to send a message to all players in a game using a wait group
func (g *IntegratedConnectionManager) Broadcast(gameId uuid.UUID, message []byte) {
	g.lock.Lock()
	defer g.lock.Unlock()

	game, found := g.GameConnectionMap[gameId]
	if !found {
		logger.Logger.Error("Cannot find game",
			"gameId", gameId)
		return
	}

	overallError := false
	var wg sync.WaitGroup
	wg.Add(len(game.playerConnectionMap))

	for playerId, conn := range game.playerConnectionMap {
		go func(conn *WsConnection, playerId uuid.UUID) {
			defer wg.Done()
			err := conn.Send(message)
			if err != nil {
				logger.Logger.Error("Cannot send a message to a player",
					"playerId", playerId)
				overallError = true

				go g.UnregisterConnection(gameId, playerId)
			}
		}(conn, playerId)
	}

	wg.Wait()
	if overallError {
		logger.Logger.Error("There was an error sending a message to a player during a broadcast operation")
	}
}

func (g *IntegratedConnectionManager) RemoveGame(gameId uuid.UUID) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	connections, found := g.GameConnectionMap[gameId]
	if !found {
		return errors.New("Cannot find the game")
	}

	var wg sync.WaitGroup
	logger.Logger.Info("Removing players from game", "gameId", gameId)
	for id, conn := range connections.playerConnectionMap {
		wg.Add(1)

		go func(id uuid.UUID, conn *WsConnection) {
			defer wg.Done()
			logger.Logger.Info("Removing player from game",
				"playerId", id,
				"gameId", gameId)
			conn.Close()
		}(id, conn)
	}
	wg.Wait()

	delete(g.GameConnectionMap, gameId)
	return nil
}

func (g *IntegratedConnectionManager) RemovePlayer(gameId, playerId uuid.UUID) error {
	res, err := gameRepo.Repo.PlayerLeaveGame(gameId, playerId)
	if err != nil {
		logger.Logger.Error("Cannot remove player from game", "err", err)
		return err
	}

	g.UnregisterConnection(gameId, playerId)

	if res.PlayersLeft == 0 {
		g.RemoveGame(gameId)
	}

	var nilUuid uuid.UUID
	if res.NewGameOwner != nilUuid {
		msg := RpcNewOwnerMsg{Id: res.NewGameOwner}
		message, err := EncodeRpcMessage(msg)
		if err != nil {
			logger.Logger.Error("Cannot encode the message",
				"err", err)
			return err
		}

		go g.Broadcast(gameId, message)
	}

	msg := RpcOnPlayerLeaveMsg{Id: playerId,
		Reason: "Player choice",
	}
	message, err := EncodeRpcMessage(msg)
	if err != nil {
		logger.Logger.Error("Cannot encode the message",
			"err", err)
		return err
	}

	go g.Broadcast(gameId, message)
	return nil
}
