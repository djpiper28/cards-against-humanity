package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
}

func (s *ServerTestSuite) SetupSuite() {
	t := s.T()

	go Start()
	t.Log("Sleeping whils the server starts")
	time.Sleep(time.Second)
	resp, err := http.Get(HttpBaseUrl + "/healthcheck")
	assert.Nil(t, err, "There should not be an error on the started server", err)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Healthcheck should work")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.Equal(t, `{"healthy":true}`, string(body), "Should return healthy")

	// Initial state checks
	s.BeforeGetGamesNotFullEmpty()
	s.BeforeInitialGameCreateTest()
}

func TestServerStart(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) BeforeGetGamesNotFullEmpty() {
	t := s.T()

	resp, err := http.Get(HttpBaseUrl + "/games/notFull")
	assert.Nil(t, err, "There should not be an error getting the games")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.Equal(t, string(body), "[]", "Should be an empty array")
}

func (s *ServerTestSuite) BeforeInitialGameCreateTest() {
	t := s.T()
	name := "Dave"

	gid, pid, err := GameRepo.CreateGame(gameLogic.DefaultGameSettings(), name)
	assert.Nil(t, err, "Should be able to make a game")
	assert.NotEmpty(t, gid)
	assert.NotEmpty(t, pid)

	resp, err := http.Get(HttpBaseUrl + "/games/notFull")
	assert.Nil(t, err, "There should not be an error getting the games")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")

	var games []gameLogic.GameInfo
	err = json.Unmarshal(body, &games)
	assert.Nil(t, err, "There should not be an error")
	assert.Len(t, games, 1, "There should be one game")

	assert.Equal(t, games[0].Id, gid, "Game Id should match")
	assert.Equal(t, games[0].PlayerCount, 1, "Should only be one player")
}

func (s *ServerTestSuite) TestGetMetrics() {
	t := s.T()
	t.Parallel()

	resp, err := http.Get(HttpBaseUrl + "/metrics")
	assert.Nil(t, err, "There should not be an error getting the metrics")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.NotEmpty(t, body, "Body should not be empty")
}

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

func (s *ServerTestSuite) TestJoinGameEndpoint() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	url := WsBaseUrl + "/games/join" +
		"?" + JoinGameGameIdParam + "=" + game.GameId.String() +
		"&" + JoinGamePlayerIdParam + "=" + game.PlayerId.String()

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	conn, _, err := dialer.Dial(url, nil)
	assert.Nil(t, err, "Should have connected to the ws server successfully")
	defer conn.Close()
	assert.NotNil(t, conn)

	msgType, msg, err := conn.ReadMessage()
	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	var onJoinMsg onJoinRpcMsg
	err = json.Unmarshal(msg, &onJoinMsg)

	assert.Nil(t, err, "Should be a join message")
	assert.Equal(t, game.GameId, onJoinMsg.Data.State.Id)
	assert.Len(t, onJoinMsg.Data.State.Players, 1)
	assert.Contains(t, onJoinMsg.Data.State.Players, game.PlayerId)
}

func (s *ServerTestSuite) TestGetCardPacks() {
	t := s.T()
	t.Parallel()

	resp, err := http.Get(HttpBaseUrl + "/res/packs")
	assert.Nil(t, err, "Should not get an error getting packs")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return a 200")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")

	var packs map[uuid.UUID]gameLogic.CardPack
	err = json.Unmarshal(body, &packs)
	assert.Nil(t, err, "There should not be any errors getting the card packs")
}
