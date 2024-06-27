package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

func (tgc *TestGameConnection) AddPlayer(playerName string) (GameJoinCookieJar, error) {
	jsonBody := CreatePlayerRequest{
		PlayerName: playerName,
		GameId:     tgc.GameId,
	}
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return GameJoinCookieJar{}, err
	}

	resp, err := http.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(body))
	if err != nil {
		return GameJoinCookieJar{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return GameJoinCookieJar{}, err
	}

	var createdPlayerInfo CreatePlayerResponse
	err = json.Unmarshal(data, &createdPlayerInfo)

	return GameJoinCookieJar{GameId: tgc.GameId,
		PlayerId: createdPlayerInfo.PlayerId,
		Token:    createdPlayerInfo.AuthorisationCookie}, nil
}
