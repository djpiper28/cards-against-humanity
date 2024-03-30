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

  playerConnection, foundPlayer := game.playerConnectionMap[playerId]
  if foundPlayer {
    playerConnection.Conn.Close()
  }

	game.playerConnectionMap[playerId] = connection

	go connection.ListenAndHandle(g)

	name, err := GameRepo.GetPlayerName(gameId, playerId)
	if err != nil {
		log.Printf("Cannot get the player's name: %s", err)
		name = "Error"
	}

	onPlayerJoinmsg := RpcOnPlayerJoinMsg{
		Id:   playerId,
		Name: name,
	}

	message, err := EncodeRpcMessage(onPlayerJoinmsg)
	if err != nil {
		log.Printf("Cannot encode the message: %s", err)
		return
	}

	go g.Broadcast(gameId, message)
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
		go func() {
			defer wg.Done()
			err := conn.Conn.Send(message)
			if err != nil {
				log.Printf("Cannot send a message to %s", playerId)
				overallError = true

				go g.UnregisterConnection(gameId, playerId)
			}
		}()
	}

	wg.Wait()
	if overallError {
		log.Print("There was an error sending a message to a player during a broadcast operation")
	}
}
