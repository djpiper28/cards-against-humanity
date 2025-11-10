package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func (s *ServerTestSuite) TestCreateGameEndpoint() {
	t := s.T()
	t.Parallel()

	name := "Dave"
	gs := DefaultGameSettings()

	postBody, err := json.Marshal(GameCreateRequest{Settings: gs, PlayerName: name})
	require.Nil(t, err, "Should be able to create json body")

	reader := bytes.NewReader(postBody)

	jar := &GameJoinCookieJar{}
	client := &http.Client{
		Timeout: time.Second * 10,
		Jar:     jar,
	}

	resp, err := client.Post(HttpBaseUrl+"/games/create", jsonContentType, reader)
	require.Nil(t, err, "Should be able to POST")
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Game should have been made and is ready for connecting to")

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err, "Should be able to read the body")

	var gameIds GameCreatedResp
	err = json.Unmarshal(body, &gameIds)
	require.Nil(t, err, "There should not be an error reading the game ids")
	require.NotEmpty(t, gameIds.GameId, "Game ID should be set")
	require.NotEmpty(t, gameIds.PlayerId, "Player ID should be set")
	require.NotEmpty(t, jar, "Token should be set")

	game, err := gameRepo.Repo.GetGame(gameIds.GameId)
	require.NoError(t, err)
	require.False(t, game.PlayersMap[gameIds.PlayerId].Connected)
}

func (s *ServerTestSuite) TestCommandError() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	require.NoError(t, err)

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	require.Nil(t, err, "Should be a join message")
	require.Equal(t, client.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	client.Write([]byte(`{"type":1,"data":{"command":"start"}}`))
	_, msg, err = client.Read()
	require.Nil(t, err, "Should be able to read the message")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcCommandErrorMsg](msg)
	require.Nil(t, err, "Should be a command error message")
	require.NotEmpty(t, rpcMsg.Reason)
}

func (s *ServerTestSuite) TestJoinGameEndpoint() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	require.Nil(t, err, "Should have connected to the ws server successfully")

	// First message should be the player join broadcast
	msgType, msg, err := client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	gameRepoGame, err := gameRepo.Repo.GetGame(client.GameId)
	require.NoError(t, err)
	require.True(t, gameRepoGame.PlayersMap[client.PlayerId].Connected)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Nil(t, err)
	require.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)

	require.Nil(t, err, "Should be a join message")
	require.Equal(t, client.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})
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
	require.NotNil(t, err, "Should have connected to the ws server successfully")
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
	require.NotNil(t, err, "Should have connected to the ws server successfully")
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
	require.NotNil(t, err, "Should have connected to the ws server successfully")
}
