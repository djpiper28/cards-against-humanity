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
	game.playerConnectionMap[playerId] = connection
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
}
