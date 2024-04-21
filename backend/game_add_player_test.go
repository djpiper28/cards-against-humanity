package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func (s *ServerTestSuite) TestCreatePlayerValid() {
	t := s.T()
	t.Parallel()

	details := s.CreateDefaultGame()
	assert.NotEmpty(t, s.CreatePlayer(details.gameId, "Bob", ""))
}

func (s *ServerTestSuite) TestCreatePlayerInvalidBodyFails() {
	t := s.T()
	t.Parallel()

	resp, err := http.Post(HttpBaseUrl+"/games/join", jsonContentType, strings.NewReader("aaaaaaaa"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func (s *ServerTestSuite) TestCreatePlayerDuplicateNameFails() {
	t := s.T()
	t.Parallel()

	const name = "Bob"
	details := s.CreateDefaultGame()
	assert.NotEmpty(t, s.CreatePlayer(details.gameId, name, ""))

	jsonBody := CreatePlayerRequest{
		PlayerName: name,
		GameId:     details.gameId,
	}
	body, err := json.Marshal(jsonBody)
	assert.Nil(t, err)

	resp, err := http.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(body))
	assert.Nil(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func (s *ServerTestSuite) TestCreatePlayerInvalidPasswordFails() {
	t := s.T()
	t.Parallel()

	const name = "Bob"
	details := s.CreateDefaultGame()

	jsonBody := CreatePlayerRequest{
		PlayerName: name,
		GameId:     details.gameId,
		Password:   "wrong password",
	}
	body, err := json.Marshal(jsonBody)
	assert.Nil(t, err)

	resp, err := http.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(body))
	assert.Nil(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func (s *ServerTestSuite) TestCreatePlayerGameFullFails() {
	t := s.T()
	t.Parallel()

	details := s.CreateDefaultGame()

	// A player has already joined on line 275
	var i uint
	for i = 1; i < details.maxPlayers; i++ {
		assert.NotEmpty(t, s.CreatePlayer(details.gameId, fmt.Sprintf("Player #%d", i), ""))
	}

	jsonBody := CreatePlayerRequest{
		PlayerName: "BadBay269",
		GameId:     details.gameId,
	}
	body, err := json.Marshal(jsonBody)
	assert.Nil(t, err)

	resp, err := http.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(body))
	assert.Nil(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func (s *ServerTestSuite) TestCreateJoinAndLeaveMessagesAreSent() {
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

	var onPlayerJoinMsg onPlayerJoinMsg
	err = json.Unmarshal(msg, &onPlayerJoinMsg)
	assert.Equal(t, game.PlayerId, onPlayerJoinMsg.Data.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = conn.ReadMessage()
	assert.Equal(t, msgType, websocket.TextMessage)

	var onJoinMsg onJoinRpcMsg
	err = json.Unmarshal(msg, &onJoinMsg)

	assert.Nil(t, err, "Should be a join message")
	assert.Equal(t, game.GameId, onJoinMsg.Data.State.Id)

	// Check that player create is sent
	// Create the player
	createPlayerReq := CreatePlayerRequest{
		PlayerName: "Bob",
		GameId:     game.GameId,
	}
	createPlayerReqBody, err := json.Marshal(createPlayerReq)
	assert.Nil(t, err)

	resp, err := http.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(createPlayerReqBody))
	assert.Nil(t, err)

	playerIdBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	newPlayerId, err := uuid.ParseBytes(playerIdBytes)
	assert.Nil(t, err)

	// Read the create message
	msgType, msg, err = conn.ReadMessage()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	var onCreateMsg onPlayerCreateMsg
	err = json.Unmarshal(msg, &onCreateMsg)
	assert.Nil(t, err)

	assert.Equal(t, newPlayerId, onCreateMsg.Data.Id)
	assert.Equal(t, createPlayerReq.PlayerName, onCreateMsg.Data.Name)

	// Check that player join is sent
	cookies = GameJoinParams{GameId: game.GameId, PlayerId: newPlayerId, Password: ""}

	dialerPlayer := websocket.DefaultDialer
	dialerPlayer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	newPlayerConn, _, err := dialerPlayer.Dial(url, cookies.Headers())
	assert.Nil(t, err, "Should have connected to the ws server successfully")
	assert.NotNil(t, conn)

	// Check that the conn message is sent
	msgType, msg, err = conn.ReadMessage()
	assert.Nil(t, err, "Should be able to read the message")
	assert.Equal(t, msgType, websocket.TextMessage)

	err = json.Unmarshal(msg, &onPlayerJoinMsg)
	assert.Nil(t, err)

	assert.Equal(t, newPlayerId, onPlayerJoinMsg.Data.Id)
	assert.Equal(t, createPlayerReq.PlayerName, onPlayerJoinMsg.Data.Name)

	// Check taht the leave message is sent
	newPlayerConn.Close()

	msgType, msg, err = conn.ReadMessage()
	assert.Nil(t, err, "Should be able to read the message")
	assert.Equal(t, msgType, websocket.TextMessage)

	var onDisconnectMsg onPlayerDisconnectMsg
	err = json.Unmarshal(msg, &onDisconnectMsg)
	assert.Nil(t, err)

	assert.Equal(t, newPlayerId, onDisconnectMsg.Data.Id)
}
