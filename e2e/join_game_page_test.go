package main

import (
	"log"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *WithServicesSuite) TestJoinGameRedirectsOnEmptyGameId() {
	t := s.T()
	browser := GetBrowser()
	defer browser.Close()

	joinGame := NewJoinGamePage(browser, "")
	time.Sleep(Timeout)
	assert.Equal(t, GetBasePage(), joinGame.Page.MustInfo().URL)
}

func (s *WithServicesSuite) TestJoinGameRedirectsOnEmptyPlayerId() {
	t := s.T()

	gameId := "testing123"
	browser := GetBrowser()
	defer browser.Close()

	joinGame := NewJoinGamePage(browser, gameId)
	time.Sleep(Timeout)
	assert.Equal(t, GetBasePage()+"game/playerJoin?gameId="+gameId, joinGame.Page.MustInfo().URL)

	_, err := UpgradeFromJoinPage(joinGame)
	assert.Nil(t, err)
}

func (s *WithServicesSuite) TestGamesShowTheSameInitialSettings() {
	browser := GetBrowser()
	defer browser.Close()

	createPage := NewCreateGamePage(browser)
	assert.NotNil(s.T(), createPage, "Page should render and not be nil")
	createPage.InsertDefaultValidSettings()

	log.Print("Creating game")
	createPage.CreateGame()

	time.Sleep(Timeout)
	assert.True(s.T(), strings.Contains(createPage.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	adminLobbyPage := JoinGamePage{Page: createPage.Page}

	assert.True(s.T(), adminLobbyPage.InLobby())
	assert.Equal(s.T(), adminLobbyPage.AdminMaxPlayers().MustText(), "6")
	assert.Equal(s.T(), adminLobbyPage.AdminPointsToPlayTo().MustText(), "10")
	assert.Equal(s.T(), adminLobbyPage.AdminMaxGameRounds().MustText(), "25")
	assert.Equal(s.T(), adminLobbyPage.AdminGamePasssowrd().MustText(), DefaultPassword)

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(s.T(), playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{Page: playerPage.Page}

	playerPage.Page.MustScreenshotFullPage("error.png")
	assert.True(s.T(), playerLobbyPage.InLobby())

	assert.Equal(s.T(),
		adminLobbyPage.AdminMaxPlayers().MustText(),
		playerLobbyPage.UserMaxPlayers().MustText())
	assert.Equal(s.T(),
		adminLobbyPage.AdminPointsToPlayTo().MustText(),
		playerLobbyPage.UserPlayingToPoints().MustText())
	assert.Equal(s.T(),
		adminLobbyPage.AdminMaxGameRounds().MustText(),
		playerLobbyPage.UserMaxGameRounds().MustText())
	assert.Equal(s.T(),
		adminLobbyPage.AdminGamePasssowrd().MustText(),
		playerLobbyPage.UserGamePassword().MustText())

	playerLobbyPage.Page.MustScreenshotFullPage("../wiki/player_lobby_page.png")
}
