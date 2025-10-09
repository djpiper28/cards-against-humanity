package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func (s *ServerTestSuite) TestLeaveGame() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	assert.Nil(t, err)
	assert.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	assert.NoError(t, err, "Should be a join message")
	assert.Equal(t, client.GameId, onJoinMsg.State.Id)
	assert.Len(t, onJoinMsg.State.Players, 1)
	assert.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	// Add the player who will leave
	name := "Player who will leave"
	leavingPlayerInfo, err := client.AddPlayer(name)
	assert.NoError(t, err)

	_, msg, err = client.Read()
	assert.Nil(t, err, "Should be able to read the message")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	assert.NoError(t, err)
	assert.Equal(t, leavingPlayerInfo.PlayerId, rpcMsg.Id)
	assert.Equal(t, name, rpcMsg.Name)

	// Leave the game
	leaveClient := http.Client{
		Jar: &leavingPlayerInfo,
	}

	leaveUrl, err := url.Parse(HttpBaseUrl + "/games/leave")
	require.NoError(t, err)

	req := &http.Request{
		URL:    leaveUrl,
		Method: http.MethodDelete,
	}
	resp, err := leaveClient.Do(req)
	assert.Nil(t, err, "There should not be an error getting the games")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	onLeaveMsg, err := network.DecodeAs[network.RpcOnPlayerLeaveMsg](msg)
	assert.NoError(t, err, "Should be a join message")
	assert.Equal(t, client.GameId, onLeaveMsg.Id)
}
