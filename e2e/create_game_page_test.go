package main

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateGamePageRender(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	page := NewCreateGamePage(browser)
	assert.NotNil(t, page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage(WikiUriBase + "create_game.png")
}

func TestCreateGamePageDefaultInput(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	page := NewCreateGamePage(browser)
	assert.NotNil(t, page, "Page should render and not be nil")
	page.InsertDefaultValidSettings()
	page.Page.MustScreenshotFullPage(WikiUriBase + "create_game_default_input.png")

	log.Print("Creating game")
	page.CreateGame()

	time.Sleep(Timeout)
	assert.True(t, strings.Contains(page.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	lobbyPage := JoinGamePage{PlayerJoinGame{page.Page}}

	assert.True(t, lobbyPage.InLobbyAdmin())
	assert.Equal(t, lobbyPage.AdminMaxPlayers().MustText(), "6")
	assert.Equal(t, lobbyPage.AdminPointsToPlayTo().MustText(), "10")
	assert.Equal(t, lobbyPage.AdminMaxGameRounds().MustText(), "25")
	assert.Equal(t, lobbyPage.AdminGamePasssowrd().MustText(), "poop")
}
