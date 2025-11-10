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
	"github.com/stretchr/testify/require"
)

func (s *ServerTestSuite) BeforeGetGamesNotFullEmpty() {
	t := s.T()

	resp, err := http.Get(HttpBaseUrl + "/games/notFull")
	require.Nil(t, err, "There should not be an error getting the games")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err, "Should be able to read the body")
	require.Equal(t, string(body), "[]", "Should be an empty array")
}

func (s *ServerTestSuite) BeforeInitialGameCreateTest() {
	t := s.T()
	name := "Dave"

	gid, pid, err := gameRepo.Repo.CreateGame(gameLogic.DefaultGameSettings(), name)
	require.Nil(t, err, "Should be able to make a game")
	require.NotEmpty(t, gid)
	require.NotEmpty(t, pid)

	resp, err := http.Get(HttpBaseUrl + "/games/notFull")
	require.Nil(t, err, "There should not be an error getting the games")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err, "Should be able to read the body")

	var games []gameLogic.GameInfo
	err = json.Unmarshal(body, &games)
	require.Nil(t, err, "There should not be an error")
	require.Len(t, games, 1, "There should be one game")

	require.Equal(t, games[0].Id, gid, "Game Id should match")
	require.Equal(t, games[0].PlayerCount, 1, "Should only be one player")
}

func (s *ServerTestSuite) TestGetMetrics() {
	t := s.T()
	t.Parallel()

	resp, err := http.Get(HttpBaseUrl + "/metrics")
	require.Nil(t, err, "There should not be an error getting the metrics")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err, "Should be able to read the body")
	require.NotEmpty(t, body, "Body should not be empty")
}

func (s *ServerTestSuite) TestGetCardPacks() {
	t := s.T()
	t.Parallel()

	resp, err := http.Get(HttpBaseUrl + "/res/packs")
	require.Nil(t, err, "Should not get an error getting packs")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Should return a 200")

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err, "Should be able to read the body")

	var packs map[uuid.UUID]gameLogic.CardPack
	err = json.Unmarshal(body, &packs)
	require.Nil(t, err, "There should not be any errors getting the card packs")
}

func (s *ServerTestSuite) TestLeaveGame() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Nil(t, err)
	require.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	require.NoError(t, err, "Should be a join message")
	require.Equal(t, client.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	// Add the player who will leave
	name := "Player who will leave"
	leavingPlayerInfo, err := client.AddPlayer(name)
	require.NoError(t, err)

	_, msg, err = client.Read()
	require.Nil(t, err, "Should be able to read the message")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.NoError(t, err)
	require.Equal(t, leavingPlayerInfo.PlayerId, rpcMsg.Id)
	require.Equal(t, name, rpcMsg.Name)

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
	require.Nil(t, err, "There should not be an error getting the games")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	onLeaveMsg, err := network.DecodeAs[network.RpcOnPlayerLeaveMsg](msg)
	require.NoError(t, err, "Should be a join message")
	require.Equal(t, client.GameId, onLeaveMsg.Id)
}
