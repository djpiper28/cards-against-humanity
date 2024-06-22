package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type TestGameConnection struct {
	GameId, PlayerId uuid.UUID
	Connection       *websocket.Conn
}

func (tgc *TestGameConnection) Close() {
	tgc.Connection.Close()
}

const WsMessageType = websocket.TextMessage

func (tgc *TestGameConnection) Read() (int, []byte, error) {
	return tgc.Connection.ReadMessage()
}

func (tgc *TestGameConnection) Write(msg []byte) error {
	return tgc.Connection.WriteMessage(WsMessageType, msg)
}

func NewTestGameConnection() (*TestGameConnection, error) {
	game, err := createTestGame_2()
	if err != nil {
		return nil, err
	}

	const url = WsBaseUrl + "/games/join"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	conn, _, err := dialer.Dial(url, game.Jar.Headers())
	if err != nil {
		log.Printf("Error dialing server: %v", err)
		return nil, err
	}

	return &TestGameConnection{GameId: game.Ids.GameId,
		PlayerId:   game.Ids.PlayerId,
		Connection: conn,
	}, nil
}
