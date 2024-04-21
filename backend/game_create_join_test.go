package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func (s *ServerTestSuite) TestCreateGameEndpoint() {
	t := s.T()
	t.Parallel()

	name := "Dave"
	gs := DefaultGameSettings()

	postBody, err := json.Marshal(GameCreateRequest{Settings: gs, PlayerName: name})
	assert.Nil(t, err, "Should be able to create json body")

	reader := bytes.NewReader(postBody)

	resp, err := http.Post(HttpBaseUrl+"/games/create", jsonContentType, reader)
	assert.Nil(t, err, "Should be able to POST")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Game should have been made and is ready for connecting to")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")

	var gameIds GameCreatedResp
	err = json.Unmarshal(body, &gameIds)
	assert.Nil(t, err, "There should not be an error reading the game ids")
	assert.NotEmpty(t, gameIds.GameId, "Game ID should be set")
	assert.NotEmpty(t, gameIds.PlayerId, "Player ID should be set")
}

func (s *ServerTestSuite) TestCommandError() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	url := WsBaseUrl + "/games/join"
	cookies := GameJoinParams{GameId: game.GameId, PlayerId: game.PlayerId, Password: ""}

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	conn, _, err := dialer.Dial(url, cookies.Headers())
	assert.Nil(t, err, "Should have connected to the ws server successfully")
	defer conn.Close()
	assert.NotNil(t, conn)

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := conn.ReadMessage()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	var onPlayerJoinMsg onPlayerJoinMsg
	err = json.Unmarshal(msg, &onPlayerJoinMsg)
	assert.Nil(t, err)
	assert.Equal(t, game.PlayerId, onPlayerJoinMsg.Data.Id, "The current user should have joined the game")

	// Second message should be the state

	msgType, msg, err = conn.ReadMessage()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	var onJoinMsg onJoinRpcMsg
	err = json.Unmarshal(msg, &onJoinMsg)

	assert.Nil(t, err, "Should be a join message")
	assert.Equal(t, game.GameId, onJoinMsg.Data.State.Id)
	assert.Len(t, onJoinMsg.Data.State.Players, 1)
	assert.Contains(t, onJoinMsg.Data.State.Players, gameLogic.Player{
		Id:   game.PlayerId,
		Name: "Dave"})

	conn.WriteMessage(websocket.TextMessage, []byte(`{"type":1,"data":{"command":"start"}}`))
	_, msg, err = conn.ReadMessage()
	assert.Nil(t, err, "Should be able to read the message")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")

	var rpcMsg onCommandError
	err = json.Unmarshal(msg, &rpcMsg)

	assert.Nil(t, err, "Should be a command error message")
	assert.NotEmpty(t, rpcMsg.Data.Reason)
}

func (s *ServerTestSuite) TestJoinGameEndpoint() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	url := WsBaseUrl + "/games/join"
	cookies := GameJoinParams{GameId: game.GameId, PlayerId: game.PlayerId, Password: ""}

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	conn, _, err := dialer.Dial(url, cookies.Headers())
	assert.Nil(t, err, "Should have connected to the ws server successfully")
	defer conn.Close()
	assert.NotNil(t, conn)

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := conn.ReadMessage()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	var onPlayerJoinMsg onPlayerJoinMsg
	err = json.Unmarshal(msg, &onPlayerJoinMsg)
	assert.Nil(t, err)
	assert.Equal(t, game.PlayerId, onPlayerJoinMsg.Data.Id, "The current user should have joined the game")

	// Second message should be the state

	msgType, msg, err = conn.ReadMessage()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	var onJoinMsg onJoinRpcMsg
	err = json.Unmarshal(msg, &onJoinMsg)

	assert.Nil(t, err, "Should be a join message")
	assert.Equal(t, game.GameId, onJoinMsg.Data.State.Id)
	assert.Len(t, onJoinMsg.Data.State.Players, 1)
	assert.Contains(t, onJoinMsg.Data.State.Players, gameLogic.Player{
		Id:   game.PlayerId,
		Name: "Dave"})
}
