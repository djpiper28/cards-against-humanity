package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/network"
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

	gid, pid, err := network.GameRepo.CreateGame(gameLogic.DefaultGameSettings(), name)
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

// A cookie jar and cookie header implementation for the ws dailer and http clients
type GameJoinParams struct {
	GameId, PlayerId uuid.UUID
	Password         string
}

func (g *GameJoinParams) Headers() http.Header {
	headers := make(http.Header)
	headers["Cookie"] = []string{fmt.Sprintf("%s=%s; %s=%s; %s=%s", JoinGamePlayerIdParam, g.PlayerId, JoinGameGameIdParam, g.GameId, PasswordParam, g.Password)}
	return headers
}

func (g *GameJoinParams) SetCookies(u *url.URL, cookies []*http.Cookie) {
	log.Fatal("Setting cookies is not desired for this API")
}

func (g *GameJoinParams) Cookies() []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, &http.Cookie{Name: JoinGamePlayerIdParam, Value: g.PlayerId.String()})
	cookies = append(cookies, &http.Cookie{Name: JoinGameGameIdParam, Value: g.GameId.String()})
	cookies = append(cookies, &http.Cookie{Name: PasswordParam, Value: g.Password})
	return cookies
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

	msgType, msg, err := conn.ReadMessage()
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

type testGameInfo struct {
	gameId, playerId uuid.UUID
	maxPlayers       uint
	password         string
}

func (s *ServerTestSuite) createDefaultGame() testGameInfo {
	t := s.T()

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

	return testGameInfo{gameId: gameIds.GameId, playerId: gameIds.PlayerId, maxPlayers: gs.MaxPlayers, password: gs.Password}
}

func (s *ServerTestSuite) createPlayer(gameId uuid.UUID, name string) uuid.UUID {
	t := s.T()

	jsonBody := CreatePlayerRequest{
		PlayerName: name,
		GameId:     gameId,
	}
	body, err := json.Marshal(jsonBody)
	assert.Nil(t, err)

	resp, err := http.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(body))
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	playerId, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	assert.NotEmpty(t, playerId)

	playerIdAsUUID, err := uuid.Parse(string(playerId))
	assert.Nil(t, err)
	return playerIdAsUUID
}

func (s *ServerTestSuite) TestCreatePlayerValid() {
	t := s.T()
	t.Parallel()

	details := s.createDefaultGame()
	assert.NotEmpty(t, s.createPlayer(details.gameId, "Bob"))
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
	details := s.createDefaultGame()
	assert.NotEmpty(t, s.createPlayer(details.gameId, name))

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

func (s *ServerTestSuite) TestCreatePlayerGameFullFails() {
	t := s.T()
	t.Parallel()

	details := s.createDefaultGame()

	// A player has already joined on line 275
	var i uint
	for i = 0; i < details.maxPlayers-1; i++ {
		assert.NotEmpty(t, s.createPlayer(details.gameId, fmt.Sprintf("Player #%d", i)))
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
