package main

import (
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameConnection struct {
}

// Manages all of the connections
var globalConnectionManager = &GlobalConnectionManager{GameConnectionMap: make(map[uuid.UUID]*GameConnection)}

type GlobalConnectionManager struct {
	// Maps a game ID to the game connection pool
	GameConnectionMap map[uuid.UUID]*GameConnection
	lock              sync.Mutex
}

type ConnectionManager interface {
	RegisterConnection(gameId, playerId uuid.UUID, connection *WsConnection)
	NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection
	Close(gameId, playerId uuid.UUID) error
}

func (g *GlobalConnectionManager) Close(gameId, playerId uuid.UUID) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	_, found := g.GameConnectionMap[gameId]
	if !found {
		return errors.New("Cannot close a player connection from a game that has been closed")
	}

	//TODO: close the game
	return nil
}

func (g *GlobalConnectionManager) RegisterConnection(gameId, playerId uuid.UUID, connection *WsConnection) {
	g.lock.Lock()
	defer g.lock.Unlock()

	log.Printf("Registering player %s to game %s", playerId, gameId)
	_, found := g.GameConnectionMap[gameId]
	if !found {
		log.Printf("Registered game %s", gameId)
	}

	//TODO: make a game connection snd start listening to the websocket
}
