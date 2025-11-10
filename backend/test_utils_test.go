package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const BaseUrl = "localhost:8080"
const HttpBaseUrl = "http://" + BaseUrl
const WsBaseUrl = "ws://" + BaseUrl

type ServerTestSuite struct {
	suite.Suite
}

func (s *ServerTestSuite) SetupSuite() {
	t := s.T()

	go Start()
	t.Log("Sleeping while the server starts")
	time.Sleep(time.Second)
	resp, err := http.Get(HttpBaseUrl + "/healthcheck")
	require.Nil(t, err, "There should not be an error on the started server", err)
	require.Equal(t, resp.StatusCode, http.StatusOK, "Healthcheck should work")

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err, "Should be able to read the body")
	require.Equal(t, `{"healthy":true}`, string(body), "Should return healthy")

	// Initial state checks
	s.BeforeGetGamesNotFullEmpty()
	s.BeforeInitialGameCreateTest()
}

func TestServerStart(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServerTestSuite))
}

type TestGameData struct {
	Ids GameCreatedResp
	Jar *GameJoinCookieJar
}

func createTestGame_2() (TestGameData, error) {
	name := "Dave"
	gs := DefaultGameSettings()

	postBody, err := json.Marshal(GameCreateRequest{Settings: gs, PlayerName: name})
	if err != nil {
		return TestGameData{}, err
	}

	reader := bytes.NewReader(postBody)

	jar := &GameJoinCookieJar{}
	client := &http.Client{
		Timeout: time.Second * 10,
		Jar:     jar,
	}
	resp, err := client.Post(HttpBaseUrl+"/games/create", jsonContentType, reader)
	if err != nil {
		return TestGameData{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return TestGameData{}, fmt.Errorf("Game should have been made and is ready for connecting")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestGameData{}, err
	}

	var gameIds GameCreatedResp
	err = json.Unmarshal(body, &gameIds)
	if err != nil {
		return TestGameData{}, err
	}

	if gameIds.PlayerId == uuid.Nil {
		return TestGameData{}, fmt.Errorf("Player ID should be set")
	}

	if gameIds.GameId == uuid.Nil {
		return TestGameData{}, fmt.Errorf("Game ID should be set")
	}

	if jar.Token == "" {
		return TestGameData{}, fmt.Errorf("Authorisation cookie was not set")
	}

	jar.GameId = gameIds.GameId
	jar.PlayerId = gameIds.PlayerId
	jar.Password = gs.Password
	return TestGameData{Ids: gameIds, Jar: jar}, nil
}

// This game has no password
func createTestGame(t *testing.T) TestGameData {
	data, err := createTestGame_2()
	require.Nil(t, err)
	return data
}

const jsonContentType = "application/json"

func DefaultGameSettings() GameCreateSettings {
	settings := gameLogic.DefaultGameSettings()
	return GameCreateSettings{MaxRounds: settings.MaxRounds,
		MaxPlayers:      settings.MaxPlayers,
		PlayingToPoints: settings.PlayingToPoints,
		CardPacks:       settings.CardPacks}
}

// A cookie jar and cookie header implementation for the ws dailer and http clients
type GameJoinCookieJar struct {
	GameId, PlayerId uuid.UUID
	Password         string
	Token            string
}

func (jar *GameJoinCookieJar) Headers() http.Header {
	headers := make(http.Header)
	headers["Cookie"] = []string{fmt.Sprintf("%s=%s; %s=%s; %s=%s; %s=%s;",
		JoinGamePlayerIdParam, jar.PlayerId,
		JoinGameGameIdParam, jar.GameId,
		PasswordParam, jar.Password,
		AuthorizationCookie, jar.Token)}
	return headers
}

func (jar *GameJoinCookieJar) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	for _, c := range cookies {
		if c.Name == AuthorizationCookie {
			jar.Token = c.Value
		} else {
			logger.Logger.Info("Ignoring", "cookie", c.Name)
		}
	}
}

func (jar *GameJoinCookieJar) Cookies(_ *url.URL) []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, &http.Cookie{Name: JoinGamePlayerIdParam, Value: jar.PlayerId.String()})
	cookies = append(cookies, &http.Cookie{Name: JoinGameGameIdParam, Value: jar.GameId.String()})
	cookies = append(cookies, &http.Cookie{Name: PasswordParam, Value: jar.Password})
	return cookies
}

func Test_GameJoinCookiesIsCookieJar(t *testing.T) {
	t.Parallel()

	var jar http.CookieJar
	jar = &GameJoinCookieJar{}
	require.NotNil(t, jar, "Should be able to create a cookie jar")
}

type TestGameInfo struct {
	gameId, playerId uuid.UUID
	maxPlayers       uint
	password         string
}

func (s *ServerTestSuite) CreateDefaultGame() TestGameInfo {
	t := s.T()
	t.Helper()

	name := "Dave"
	gs := DefaultGameSettings()

	postBody, err := json.Marshal(GameCreateRequest{Settings: gs, PlayerName: name})
	require.Nil(t, err, "Should be able to create json body")

	reader := bytes.NewReader(postBody)

	resp, err := http.Post(HttpBaseUrl+"/games/create", jsonContentType, reader)
	require.Nil(t, err, "Should be able to POST")
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Game should have been made and is ready for connecting to")

	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err, "Should be able to read the body")

	var gameIds GameCreatedResp
	err = json.Unmarshal(body, &gameIds)
	require.Nil(t, err, "There should not be an error reading the game ids")
	require.NotEmpty(t, gameIds.GameId, "Game ID should be set")
	require.NotEmpty(t, gameIds.PlayerId, "Player ID should be set")

	return TestGameInfo{gameId: gameIds.GameId, playerId: gameIds.PlayerId, maxPlayers: gs.MaxPlayers, password: gs.Password}
}

func (s *ServerTestSuite) CreatePlayer(gameId uuid.UUID, name, password string) TestGameData {
	t := s.T()
	t.Helper()

	jsonBody := CreatePlayerRequest{
		PlayerName: name,
		GameId:     gameId,
	}
	body, err := json.Marshal(jsonBody)
	require.Nil(t, err)

	client := http.Client{
		Jar: &GameJoinCookieJar{GameId: gameId, Password: password},
	}

	resp, err := client.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(body))
	require.Nil(t, err)

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var create CreatePlayerResponse
	err = json.Unmarshal(respBody, &create)
	require.Nil(t, err)

	require.NotEmpty(t, create.PlayerId)
	return TestGameData{Ids: GameCreatedResp{GameId: gameId, PlayerId: create.PlayerId}, Jar: &GameJoinCookieJar{GameId: gameId, PlayerId: create.PlayerId, Password: password}}
}

func (s *ServerTestSuite) ReadCreateJoinMessages(t *testing.T, client *TestGameConnection, pid uuid.UUID) {
	t.Helper()

	create := false
	join := false

	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for tries := range 5 {
		var rpcMessage network.RpcMessageBody
		msgType, msg, err := client.Read()
		require.Equal(t, msgType, websocket.TextMessage)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(msg, &rpcMessage))

		switch rpcMessage.Type {
		case network.MsgOnPlayerCreate:
			onPlayerCreateMsg, err := network.DecodeAs[network.RpcOnPlayerCreateMsg](msg)
			require.Nil(t, err)
			require.Equal(t, pid, onPlayerCreateMsg.Id, "The current user should have joined the game")
			t.Log("MsgOnPlayerCreate")

			require.False(t, create)
			create = true
		case network.MsgOnPlayerJoin:
			onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
			require.Nil(t, err)
			require.Equal(t, pid, onPlayerJoinMsg.Id, "The current user should have joined the game")
			t.Log("MsgOnPlayerJoin")

			require.False(t, join)
			join = true
		case network.MsgPing:
			pong, err := network.EncodeRpcMessage(network.RpcPingMsg{})
			require.NoError(t, err)
			require.NoError(t, client.Write(pong))

			t.Log("Ping")
			continue
		case network.MsgCommandError:
			require.FailNowf(t, "MsgCommandError", "try %d msg %s", tries, string(msg))
		default:
			require.FailNowf(t, "Cannot parse message", "try %d msg %s", tries, rpcMessage.Type)
		}

		if create && join {
			return
		}
	}

	require.FailNowf(t, "Cannot read create and join messages", "create: %v, join: %v", create, join)
}

// Ignores pings that are sent
func ReadMessage[T network.RpcMessage](s *ServerTestSuite, t *testing.T, client *TestGameConnection) T {
	t.Helper()

	type TProxy struct {
		Type network.RpcMessageType `json:"type"`
		Data T                      `json:"data"`
	}

	for tries := range 5 {
		msgType, msg, err := client.Read()
		require.Equal(t, msgType, websocket.TextMessage)
		require.NoError(t, err)

		_, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		var proxy TProxy
		err = json.Unmarshal(msg, &proxy)
		require.NoError(t, err)

		if proxy.Type == network.MsgPing {
			pong, err := network.EncodeRpcMessage(network.RpcPingMsg{})
			require.NoError(t, err)
			require.NoError(t, client.Write(pong))
			t.Log("Ping")
			continue
		}

		if proxy.Type == network.MsgCommandError {
			require.FailNowf(t, "MSG Command Error was returned", "try: %d msg: %s", tries, string(msg))
		}

		require.Equal(t, proxy.Data.Type(), proxy.Type)
		return proxy.Data
	}

	require.FailNow(t, "Max tries to read desired message exceeded")
	panic("Fail")
}
