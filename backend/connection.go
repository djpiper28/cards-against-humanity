package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const wsBufferSize = 1024

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  wsBufferSize,
	WriteBufferSize: wsBufferSize,
}

func WsUpgrade(w http.ResponseWriter, r *http.Request, playerId, gameId uuid.UUID, cm ConnectionManager) (*WsConnection, error) {
	c, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %s", err)
		return nil, err
	}

	conn := cm.NewConnection(c, playerId, gameId)
	return conn, nil
}

type GameMessage struct {
	Message  string
	GameId   uuid.UUID
	PlayerId uuid.UUID
}

type WsConnection struct {
	Conn         NetworkConnection
	PlayerId     uuid.UUID
	JoinTime     time.Time
	LastPingTime time.Time
}

func (gcm *WsConnection) Close() {
  gcm.Conn.Close()
}

func (gcm *GlobalConnectionManager) NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection {
	c := &WsConnection{Conn: &WebsocketConnection{Conn: conn},
		PlayerId:     playerId,
		JoinTime:     time.Now(),
		LastPingTime: time.Now(),
	}
	gcm.RegisterConnection(gameId, playerId, c)
	return c
}
