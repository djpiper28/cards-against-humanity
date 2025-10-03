package network

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
	"github.com/google/uuid"
)

type GameConnection struct {
	// Maps a player to the connection
	playerConnectionMap map[uuid.UUID]*WsConnection
}

// Manages all of the connections
var GlobalConnectionManager = New()

func gameRemovalDaemon(manager *ConnectionManager) {
	for {
		games := gameRepo.Repo.EndOldGames()
		for _, game := range games {
			err := manager.RemoveGame(game)
			if err != nil {
				logger.Logger.Error("Cannot remove game", "gid", game, "err", err)
			}
		}
		time.Sleep(time.Second)
	}
}

func New() *ConnectionManager {
	manager := &ConnectionManager{GameConnectionMap: make(map[uuid.UUID]*GameConnection)}
	go gameRemovalDaemon(manager)
	return manager
}

type ConnectionManager struct {
	// Maps a game ID to the game connection pool
	GameConnectionMap map[uuid.UUID]*GameConnection
	lock              sync.Mutex
}

func (g *ConnectionManager) Close(gameId, playerId uuid.UUID) error {
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

func (g *ConnectionManager) RegisterConnection(gameId, playerId uuid.UUID, connection *WsConnection) {
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

func (g *ConnectionManager) UnregisterConnection(gameId, playerId uuid.UUID) {
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
			"gameId", gameId,
      "err", err)
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

func (g *ConnectionManager) SendToPlayer(gameId, playerId uuid.UUID, message []byte) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	game, found := g.GameConnectionMap[gameId]
	if !found {
		logger.Logger.Error("Cannot find game", "gameId", gameId, "playerId", playerId)
		return errors.New("Cannot find the game")
	}

	conn, found := game.playerConnectionMap[playerId]
	if found {
		return conn.Send(message)
	}
	logger.Logger.Error("Cannot find player", "gameId", gameId, "playerId", playerId)
	return errors.New("Cannot find the player")
}

// Blocking call to send a message to all players in a game using a wait group
func (g *ConnectionManager) Broadcast(gameId uuid.UUID, message []byte) {
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

func (g *ConnectionManager) RemoveGame(gameId uuid.UUID) error {
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

func (g *ConnectionManager) RemovePlayer(gameId, playerId uuid.UUID) error {
	res, err := gameRepo.Repo.PlayerLeaveGame(gameId, playerId)
	if err != nil {
		logger.Logger.Error("Cannot remove player from game", "err", err)
		return err
	}

	g.UnregisterConnection(gameId, playerId)

	if res.NewGameState == gameLogic.GameStateCzarJudgingCards {
		go g.MoveToCzarJudgingPhase(gameId, res.CzarJudingPhaseInfo)
	}

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

func (g *ConnectionManager) MoveToCzarJudgingPhase(gid uuid.UUID, info gameLogic.CzarJudingPhaseInfo) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	game, found := g.GameConnectionMap[gid]
	if !found {
		return errors.New(fmt.Sprintf("Cannot find game %s", gid))
	}

	var sendErr error

	wg := sync.WaitGroup{}
	wg.Add(len(game.playerConnectionMap))
	for pid, conn := range game.playerConnectionMap {
		go func(pid uuid.UUID, conn *WsConnection) {
			defer wg.Done()
			moveToNextPhaseMsg := RpcOnCzarJudgingPhaseMsg{
				AllPlays: info.AllPlays,
				NewHand:  info.PlayerHands.Hands[pid],
			}

			broadcastMessage, err := EncodeRpcMessage(moveToNextPhaseMsg)
			if err != nil {
				sendErr = err
				logger.Logger.Error("Cannot send a message to a player",
					"playerId", pid)
			}

			err = conn.Send(broadcastMessage)
		}(pid, conn)
	}
	wg.Wait()

	if sendErr != nil {
		logger.Logger.Warn("Was not able to send Czar judging phase message to all players")
	}
	return sendErr
}

func (gcm *ConnectionManager) Exists(gameId uuid.UUID) bool {
	_, err := gameRepo.Repo.GetGame(gameId)
	return err == nil
}
