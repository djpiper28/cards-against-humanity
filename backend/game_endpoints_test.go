package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

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

	gid, pid, err := gameRepo.Repo.CreateGame(gameLogic.DefaultGameSettings(), name)
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

func (s *ServerTestSuite) TestJoinGameEndpointFailsWrongPassword() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	game.Jar.Password = "wrong password"
	url := WsBaseUrl + "/games/join"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	_, _, err := dialer.Dial(url, game.Jar.Headers())
	assert.NotNil(t, err, "Should have connected to the ws server successfully")
}

func (s *ServerTestSuite) TestJoinGameEndpointFailsPlayerNotReal() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	game.Jar.PlayerId = uuid.New()
	url := WsBaseUrl + "/games/join"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	_, _, err := dialer.Dial(url, game.Jar.Headers())
	assert.NotNil(t, err, "Should have connected to the ws server successfully")
}

func (s *ServerTestSuite) TestJoinGameEndpointFailsGameNotReal() {
	t := s.T()
	t.Parallel()

	url := WsBaseUrl + "/games/join"
	cookies := GameJoinCookieJar{GameId: uuid.New(), PlayerId: uuid.New(), Password: ""}

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	_, _, err := dialer.Dial(url, cookies.Headers())
	assert.NotNil(t, err, "Should have connected to the ws server successfully")
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
