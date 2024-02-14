package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
