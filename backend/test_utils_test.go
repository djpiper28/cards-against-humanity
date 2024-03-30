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

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const BaseUrl = "localhost:8080"
const HttpBaseUrl = "http://" + BaseUrl
const WsBaseUrl = "ws://" + BaseUrl

// This game has no password
func createTestGame(t *testing.T) GameCreatedResp {
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
	return gameIds
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

func (g *GameJoinParams) SetCookies(_ *url.URL, _ []*http.Cookie) {
	log.Fatal("Setting cookies is not desired for this API")
}

func (g *GameJoinParams) Cookies(_ *url.URL) []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, &http.Cookie{Name: JoinGamePlayerIdParam, Value: g.PlayerId.String()})
	cookies = append(cookies, &http.Cookie{Name: JoinGameGameIdParam, Value: g.GameId.String()})
	cookies = append(cookies, &http.Cookie{Name: PasswordParam, Value: g.Password})
	return cookies
}

func Test_GameJoinCookiesIsCookieJar(t *testing.T) {
  t.Parallel()

  var jar http.CookieJar
  jar = &GameJoinParams{}
  assert.NotNil(t, jar, "Should be able to create a cookie jar")
}
