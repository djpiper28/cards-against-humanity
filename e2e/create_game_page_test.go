package main

import (
	"log"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *WithServicesSuite) TestCreateGamePageRender() {
	page := NewCreateGamePage(GetBrowser())
	assert.NotNil(s.T(), page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage("../wiki/assets/create_game.png")
}

func (s *WithServicesSuite) TestCreateGamePageDefaultInput() {
	page := NewCreateGamePage(GetBrowser())
	assert.NotNil(s.T(), page, "Page should render and not be nil")
	page.InsertDefaultValidSettings()
	page.Page.MustScreenshotFullPage("../wiki/assets/create_game_default_input.png")

  log.Print("Creating game")
	page.CreateGame()

	time.Sleep(Timeout)
	assert.True(s.T(), strings.Contains(page.Page.Timeout(Timeout).MustInfo().URL, "/join?gameId="))

	lobbyPage := JoinGamePage{Page: page.Page}
	assert.True(s.T(), lobbyPage.InLobby())
}
