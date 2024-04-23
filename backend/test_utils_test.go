package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

type GameData struct {
	Ids GameCreatedResp
	Jar *GameJoinCookieJar
}

// This game has no password
func createTestGame(t *testing.T) GameData {
	name := "Dave"
	gs := DefaultGameSettings()

	postBody, err := json.Marshal(GameCreateRequest{Settings: gs, PlayerName: name})
	assert.Nil(t, err, "Should be able to create json body")

	reader := bytes.NewReader(postBody)

	jar := &GameJoinCookieJar{}
	client := &http.Client{
		Timeout: time.Second * 10,
		Jar:     jar,
	}
	resp, err := client.Post(HttpBaseUrl+"/games/create", jsonContentType, reader)
	assert.Nil(t, err, "Should be able to POST")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Game should have been made and is ready for connecting to")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")

	var gameIds GameCreatedResp
	err = json.Unmarshal(body, &gameIds)
	assert.Nil(t, err, "There should not be an error reading the game ids")

	assert.NotEmpty(t, gameIds.PlayerId, "Player ID should be set")
	assert.NotEmpty(t, gameIds.GameId, "Game ID should be set")
	assert.NotEmpty(t, jar.Token, "Authorisation cookie was not set")

	jar.GameId = gameIds.GameId
	jar.PlayerId = gameIds.PlayerId
	jar.Password = gs.Password
	return GameData{Ids: gameIds, Jar: jar}
}

const jsonContentType = "application/json"

func DefaultGameSettings() GameCreateSettings {
	settings := gameLogic.DefaultGameSettings()
	packs := make([]uuid.UUID, len(settings.CardPacks))
	for i, pack := range settings.CardPacks {
		packs[i] = pack.Id
	}
	return GameCreateSettings{MaxRounds: settings.MaxRounds, MaxPlayers: settings.MaxPlayers, PlayingToPoints: settings.PlayingToPoints, CardPacks: packs}
}

type onJoinRpcMsg struct {
	Data network.RpcOnJoinMsg `json:"data"`
}

type onPlayerJoinMsg struct {
	Data network.RpcOnPlayerJoinMsg `json:"data"`
}

type onPlayerCreateMsg struct {
	Data network.RpcOnPlayerCreateMsg `json:"data"`
}

type onPlayerDisconnectMsg struct {
	Data network.RpcOnPlayerDisconnectMsg `json:"data"`
}

type onChangeSettings struct {
	Data network.RpcChangeSettingsMsg `json:"data"`
}

type onCommandError struct {
	Data network.RpcCommandErrorMsg `json:"data"`
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
			log.Printf("Ignoring cookie %s", c.Name)
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
	assert.NotNil(t, jar, "Should be able to create a cookie jar")
}

type TestGameInfo struct {
	gameId, playerId uuid.UUID
	maxPlayers       uint
	password         string
}

func (s *ServerTestSuite) CreateDefaultGame() TestGameInfo {
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

	return TestGameInfo{gameId: gameIds.GameId, playerId: gameIds.PlayerId, maxPlayers: gs.MaxPlayers, password: gs.Password}
}

func (s *ServerTestSuite) CreatePlayer(gameId uuid.UUID, name, password string) GameData {
	t := s.T()

	jsonBody := CreatePlayerRequest{
		PlayerName: name,
		GameId:     gameId,
	}
	body, err := json.Marshal(jsonBody)
	assert.Nil(t, err)

	client := http.Client{
		Jar: &GameJoinCookieJar{GameId: gameId, Password: password},
	}

	resp, err := client.Post(HttpBaseUrl+"/games/join", jsonContentType, bytes.NewReader(body))
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	var create CreatePlayerResponse
	err = json.Unmarshal(respBody, &create)
	assert.Nil(t, err)

	assert.NotEmpty(t, create.PlayerId)
	return GameData{Ids: GameCreatedResp{GameId: gameId, PlayerId: create.PlayerId}, Jar: &GameJoinCookieJar{GameId: gameId, PlayerId: create.PlayerId, Password: password}}
}
