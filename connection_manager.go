package main

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"sync"
)

type GameConnection struct {
	PlayerConnectionMap map[uuid.UUID]*WsConnection
	broadcast           chan string
	shutdown            chan bool
	lock                sync.Mutex
}

func NewGameConnection() *GameConnection {
	gc := &GameConnection{PlayerConnectionMap: NewGameConnection().PlayerConnectionMap,
		broadcast: make(chan string),
		shutdown:  make(chan bool)}
	go gc.process()
	return gc
}

func (gc *GameConnection) process() {
	for {
		select {
		case msg := <-gc.broadcast:
			gc.lock.Lock()
			for _, conn := range gc.PlayerConnectionMap {
				conn.WsBroadcast <- msg
			}
			gc.lock.Unlock()
		default:
			log.Printf("Shut the game down")
			return
		}
	}
}

func (gc *GameConnection) Close(playerId uuid.UUID) error {
	gc.lock.Lock()
	defer gc.lock.Unlock()

	log.Printf("Player %s disconnected", playerId)
	conn, found := gc.PlayerConnectionMap[playerId]
	if !found {
		return errors.New("Cannot remove a wensocket that has been closed already")
	}

	conn.Close()
	delete(gc.PlayerConnectionMap, playerId)
	return nil
}

func (gc *GameConnection) CloseAll() {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	close(gc.broadcast)
	for _, conn := range gc.PlayerConnectionMap {
		conn.Close()
	}
}

func (gc *GameConnection) RegisterConnection(playerId uuid.UUID, conn *WsConnection) {
	gc.lock.Lock()
	defer gc.lock.Unlock()

	player, found := gc.PlayerConnectionMap[playerId]
	if found {
		log.Printf("Evicting player %s's old connection", playerId)
		player.Close()
	}

	gc.PlayerConnectionMap[playerId] = conn
}

// Manages all of the connections
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

func (g *GlobalConnectionManager) RegisterConnection(gameId, playerId uuid.UUID, connection *WsConnection) {
	g.lock.Lock()
	defer g.lock.Unlock()

	log.Printf("Registering player %s to game %s", playerId, gameId)
	gameMap, found := g.GameConnectionMap[gameId]
	if !found {
		log.Printf("Registered game %s", gameId)
		newGameMap := NewGameConnection()
		g.GameConnectionMap[gameId] = newGameMap
		gameMap = newGameMap
	}

	gameMap.RegisterConnection(playerId, connection)
}
