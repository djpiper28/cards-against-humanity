package network

import (
	"errors"
	"log"
	"sync"

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

	log.Printf("Registering player %s to game %s", playerId, gameId)
	game, found := g.GameConnectionMap[gameId]
	if !found {
		game = &GameConnection{make(map[uuid.UUID]*WsConnection)}
		g.GameConnectionMap[gameId] = game
		log.Printf("Registered game %s", gameId)
	}

	playerConnection, foundPlayer := game.playerConnectionMap[playerId]
	if foundPlayer {
		// Do not trigger a disconnect event in the internal system, just hot-swap the connection
		playerConnection.conn.Close()
	}

	game.playerConnectionMap[playerId] = connection

	err := GameRepo.ConnectPlayer(gameId, playerId)
	if err != nil {
		log.Printf("Cannot tag player %s as connected to game %s", playerId, gameId)
	}

	go connection.ListenAndHandle(g)
}

func (g *IntegratedConnectionManager) UnregisterConnection(gameId, playerId uuid.UUID) {
	g.lock.Lock()
	defer g.lock.Unlock()

	log.Printf("Unregistering player %s to game %s", playerId, gameId)
	game, found := g.GameConnectionMap[gameId]
	if found {
		delete(game.playerConnectionMap, playerId)
	} else {
		log.Printf("Cannot unregister game %s as it cannot be found", gameId)
	}

	err := GameRepo.DisconnectPlayer(gameId, playerId)
	if err != nil {
		log.Printf("Cannot tag player %s as disconnected from game %s", playerId, gameId)
	}

	onPlayerDisconnectMsg := RpcOnPlayerDisconnectMsg{
		Id: playerId,
	}

	message, err := EncodeRpcMessage(onPlayerDisconnectMsg)
	if err != nil {
		log.Printf("Cannot encode the message: %s", err)
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
		log.Printf("Cannot find game: %s", gameId)
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
				log.Printf("Cannot send a message to %s", playerId)
				overallError = true

				go g.UnregisterConnection(gameId, playerId)
			}
		}(conn, playerId)
	}

	wg.Wait()
	if overallError {
		log.Print("There was an error sending a message to a player during a broadcast operation")
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
	log.Printf("Removing players from game %s", gameId)
	for id, conn := range connections.playerConnectionMap {
		wg.Add(1)

		go func(id uuid.UUID, conn *WsConnection) {
			defer wg.Done()
			log.Printf("Removing player %s from game %s", id, gameId)
			conn.Close()
		}(id, conn)
	}
	wg.Wait()

	delete(g.GameConnectionMap, gameId)
	return nil
}

func (g *IntegratedConnectionManager) RemovePlayer(gameId, playerId uuid.UUID) error {
	res, err := GameRepo.PlayerLeaveGame(gameId, playerId)
	if err != nil {
		log.Printf("Cannot remove player from game: %s", err)
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
			log.Printf("Cannot encode the message: %s", err)
			return err
		}

		go g.Broadcast(gameId, message)
	}

	msg := RpcOnPlayerLeaveMsg{Id: playerId,
		Reason: "Player choice",
	}
	message, err := EncodeRpcMessage(msg)
	if err != nil {
		log.Printf("Cannot encode the message: %s", err)
		return err
	}

	go g.Broadcast(gameId, message)
	return nil
}
