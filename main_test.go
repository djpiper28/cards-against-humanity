package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/stretchr/testify/assert"
)

const baseUrl = "http://localhost:8080"

func TestServerStart(t *testing.T) {
	go Start()

	t.Log("Sleeping whils the server starts")
	time.Sleep(time.Second)
	resp, err := http.Get(baseUrl + "/healthcheck")
	assert.Nil(t, err, "There should not be an error on the started server", err)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Healthcheck should work")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.Equal(t, `{"healthy":true}`, string(body), "Should return healthy")
}

func TestGetMetrics(t *testing.T) {
	resp, err := http.Get(baseUrl + "/metrics")
	assert.Nil(t, err, "There should not be an error getting the metrics")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.NotEmpty(t, body, "Body should not be empty")
}

func TestGetGamesNotFullEmpty(t *testing.T) {
	resp, err := http.Get(baseUrl + "/games/notFull")
	assert.Nil(t, err, "There should not be an error getting the games")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.Equal(t, string(body), "[]", "Should be an empty array")
}

func TestGetGamesNotFullOneGame(t *testing.T) {
	name := "Dave"

	gid, pid, err := GameRepo.CreateGame(gameLogic.DefaultGameSettings(), name)
	assert.Nil(t, err, "Should be able to make a game")
	assert.NotEmpty(t, gid)
	assert.NotEmpty(t, pid)

	resp, err := http.Get(baseUrl + "/games/notFull")
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

const jsonContentType = "application/json"

func TestCreateGameEndpoint(t *testing.T) {
	name := "Dave"
	gs := gameLogic.DefaultGameSettings()

	postBody, err := json.Marshal(gameCreateSettings{Settings: *gs, PlayerName: name})
	assert.Nil(t, err, "Should be able to create json body")

	reader := bytes.NewReader(postBody)

	resp, err := http.Post(baseUrl+"/games/create", jsonContentType, reader)
	assert.Nil(t, err, "Should be able to POST")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Game should have been made and is ready for connecting to")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")

	var gameIds gameCreatedResp
	err = json.Unmarshal(body, &gameIds)
	assert.Nil(t, err, "There should not be an error reading the game ids")
	assert.NotEmpty(t, gameIds.GameId, "Game ID should be set")
	assert.NotEmpty(t, gameIds.PlayerId, "Player ID should be set")
}
