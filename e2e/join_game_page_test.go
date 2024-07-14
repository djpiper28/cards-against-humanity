package main

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
)

func TestJoinGameRedirectsOnEmptyGameId(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	joinGame := NewJoinGamePage(browser, "")
	time.Sleep(Timeout)
	assert.Equal(t, GetBasePage(), joinGame.Page.MustInfo().URL)
}

func TestJoinGameRedirectsOnEmptyPlayerId(t *testing.T) {
	t.Parallel()

	gameId := "testing123"
	browser := GetBrowser()
	defer browser.Close()

	joinGame := NewJoinGamePage(browser, gameId)
	time.Sleep(Timeout)
	assert.Equal(t, GetBasePage()+"game/playerJoin?gameId="+gameId, joinGame.Page.MustInfo().URL)

	_, err := UpgradeFromJoinPage(joinGame)
	assert.Nil(t, err)
}

func TestGamesShowTheSameInitialSettings(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	createPage := NewCreateGamePage(browser)
	assert.NotNil(t, createPage, "Page should render and not be nil")
	createPage.InsertDefaultValidSettings()

	log.Print("Creating game")
	createPage.CreateGame()

	time.Sleep(Timeout)
	assert.True(t, strings.Contains(createPage.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	adminLobbyPage := JoinGamePage{PlayerJoinGame{createPage.Page}}

	assert.True(t, adminLobbyPage.InLobbyAdmin())
	assert.True(t, adminLobbyPage.IsAdmin())
	assert.Equal(t, adminLobbyPage.AdminMaxPlayers().MustText(), "6")
	assert.Equal(t, adminLobbyPage.AdminPointsToPlayTo().MustText(), "10")
	assert.Equal(t, adminLobbyPage.AdminMaxGameRounds().MustText(), "25")
	assert.Equal(t, adminLobbyPage.AdminGamePasssowrd().MustText(), DefaultPassword)

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(t, playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{PlayerJoinGame{playerPage.Page}}

	assert.True(t, playerLobbyPage.InLobbyAdmin())

	assert.Equal(t,
		adminLobbyPage.AdminMaxPlayers().MustText(),
		playerLobbyPage.UserMaxPlayers().MustText())
	assert.Equal(t,
		adminLobbyPage.AdminPointsToPlayTo().MustText(),
		playerLobbyPage.UserPlayingToPoints().MustText())
	assert.Equal(t,
		adminLobbyPage.AdminMaxGameRounds().MustText(),
		playerLobbyPage.UserMaxGameRounds().MustText())
	assert.Equal(t,
		adminLobbyPage.AdminGamePasssowrd().MustText(),
		playerLobbyPage.UserGamePassword().MustText())

	playerLobbyPage.Page.MustScreenshotFullPage(WikiUriBase + "player_lobby_page.png")
}

func TestChangingSettingsSyncsBetweenClients(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	createPage := NewCreateGamePage(browser)
	assert.NotNil(t, createPage, "Page should render and not be nil")
	createPage.InsertDefaultValidSettings()

	log.Print("Creating game")
	createPage.CreateGame()

	time.Sleep(Timeout)
	assert.True(t, strings.Contains(createPage.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	adminLobbyPage := JoinGamePage{PlayerJoinGame{createPage.Page}}

	assert.True(t, adminLobbyPage.InLobbyAdmin())

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(t, playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{PlayerJoinGame{playerPage.Page}}

	assert.True(t, playerLobbyPage.InLobbyAdmin())

	adminLobbyPage.AdminGamePasssowrd().MustInput("Password 123")
	assert.Equal(t, adminLobbyPage.AdminGamePasssowrd().MustText(), "poopPassword 123")

	assert.True(t, adminLobbyPage.Saved())

	found := false
	for _, cookie := range browser.MustGetCookies() {
		if cookie.Name == "password" {
			assert.Equal(t, "poopPassword 123", cookie.Value)
			found = true
			break
		}
	}
	assert.True(t, found)

	found = false
	for _, cookie := range browser2.MustGetCookies() {
		if cookie.Name == "password" {
			assert.Equal(t, "poopPassword 123", cookie.Value)
			found = true
			break
		}
	}
	assert.True(t, found)

	assert.Equal(t,
		adminLobbyPage.AdminMaxPlayers().MustText(),
		playerLobbyPage.UserMaxPlayers().MustText())
	assert.Equal(t,
		adminLobbyPage.AdminPointsToPlayTo().MustText(),
		playerLobbyPage.UserPlayingToPoints().MustText())
	assert.Equal(t,
		adminLobbyPage.AdminMaxGameRounds().MustText(),
		playerLobbyPage.UserMaxGameRounds().MustText())
	assert.Equal(t,
		adminLobbyPage.AdminGamePasssowrd().MustText(),
		playerLobbyPage.UserGamePassword().MustText())

	playerLobbyPage.Page.MustScreenshotFullPage(WikiUriBase + "player_lobby_page.png")
}

func TestPlayerDisconnectReConnectWithNoSearchParamsInJoinLink(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	createPage := NewCreateGamePage(browser)
	assert.NotNil(t, createPage, "Page should render and not be nil")
	createPage.InsertDefaultValidSettings()

	log.Print("Creating game")
	createPage.CreateGame()

	time.Sleep(Timeout)
	assert.True(t, strings.Contains(createPage.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	adminLobbyPage := JoinGamePage{PlayerJoinGame{createPage.Page}}

	assert.True(t, adminLobbyPage.InLobbyAdmin())

	adminPlayerId := adminLobbyPage.PlayerId()
	assert.True(t,
		adminLobbyPage.PlayerConnected(adminPlayerId))

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(t, playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{PlayerJoinGame{playerPage.Page}}
	assert.True(t, playerLobbyPage.InLobbyAdmin())

	playerId := playerLobbyPage.PlayerId()

	time.Sleep(Timeout)
	assert.True(t, adminLobbyPage.PlayerConnected(playerId))

	// Disconnect and make sure the UI updates
	playerPage.Disconnect()

	time.Sleep(Timeout)
	assert.False(t, adminLobbyPage.PlayerConnected(playerId))
	assert.True(t, adminLobbyPage.PlayerConnected(adminPlayerId))

	adminLobbyPage.Page.MustScreenshot(WikiUriBase + "admin_player_disconnected.png")
	playerPage.Page.MustScreenshot(WikiUriBase + "player_disconnected.png")

	// Reconnect and make sure the UI updates
	playerPage.ReConnect()
	time.Sleep(Timeout)
	playerPage.Page.MustScreenshot(WikiUriBase + "player_reconnect.png")

	assert.True(t, adminLobbyPage.PlayerConnected(adminPlayerId))
	assert.True(t, adminLobbyPage.PlayerConnected(playerId))
	assert.True(t, playerPage.InLobbyPlayer())
}

func TestPlayerLeavesGame(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	createPage := NewCreateGamePage(browser)
	assert.NotNil(t, createPage, "Page should render and not be nil")
	createPage.InsertDefaultValidSettings()

	log.Print("Creating game")
	createPage.CreateGame()

	time.Sleep(Timeout)
	assert.True(t, strings.Contains(createPage.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	adminLobbyPage := JoinGamePage{PlayerJoinGame{createPage.Page}}

	assert.True(t, adminLobbyPage.InLobbyAdmin())

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(t, playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{PlayerJoinGame{playerPage.Page}}
	assert.True(t, playerLobbyPage.InLobbyAdmin())

	playerId := playerLobbyPage.PlayerId()

	playerLobbyPage.LeaveGame()
	assert.NotContains(t, adminLobbyPage.PlayersInGame(), playerId)
	assert.True(t, adminLobbyPage.IsAdmin())
}

// Genuine fail remind me to fix
func TestOwnerLeavingGameTransfersOwnership(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	createPage := NewCreateGamePage(browser)
	assert.NotNil(t, createPage, "Page should render and not be nil")
	createPage.InsertDefaultValidSettings()

	log.Print("Creating game")
	createPage.CreateGame()

	time.Sleep(Timeout)
	assert.True(t, strings.Contains(createPage.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	adminLobbyPage := JoinGamePage{PlayerJoinGame{createPage.Page}}

	assert.True(t, adminLobbyPage.InLobbyAdmin())

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(t, playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{PlayerJoinGame{playerPage.Page}}
	assert.True(t, playerLobbyPage.InLobbyAdmin())

	time.Sleep(Timeout)
	playerId := adminLobbyPage.PlayerId()
	adminLobbyPage.LeaveGame()

	time.Sleep(Timeout)
	assert.NotContains(t, playerLobbyPage.PlayersInGame(), playerId)
	assert.True(t, playerLobbyPage.IsAdmin())
}

func TestStartingGame(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	createPage := NewCreateGamePage(browser)
	assert.NotNil(t, createPage, "Page should render and not be nil")
	createPage.InsertDefaultValidSettings()

	log.Print("Creating game")
	createPage.CreateGame()

	time.Sleep(Timeout)
	assert.True(t, strings.Contains(createPage.Page.Timeout(Timeout).MustInfo().URL, "game/join?gameId="))

	adminLobbyPage := JoinGamePage{PlayerJoinGame{createPage.Page}}

	assert.True(t, adminLobbyPage.InLobbyAdmin())

	browsers := make([]*rod.Browser, 0)
	const playerCount = 3
	for i := 0; i < playerCount; i++ {
		// Connect with another client then assert that the settings remain equal
		playerBrowser := GetBrowser()
		browsers = append(browsers, playerBrowser)
		defer playerBrowser.Close()
		playerPage := NewPlayerGamePage(playerBrowser, adminLobbyPage)

		assert.True(t, playerPage.InPlayerJoinPage())
		playerPage.PlayerName(fmt.Sprintf("Player %d", i))
		playerPage.Password(DefaultPassword)

		playerPage.Join()

		assert.True(t, playerPage.InLobbyPlayer())
		screenshotError(playerPage.Page)
	}

	time.Sleep(Timeout)
	adminLobbyPage.Start()
	time.Sleep(Timeout)

	cards, err := adminLobbyPage.Cards()
	assert.NoError(t, err)
	assert.Len(t, cards, 7)

	for i := 0; i < playerCount; i++ {
		playerPage := NewPlayerGamePage(browsers[i], adminLobbyPage)
		cards, err := playerPage.Cards()
		assert.NoError(t, err)
		assert.Len(t, cards, 7)

		if i == playerCount-1 {
			assert.True(t, playerPage.IsCzar())
		}
	}
}
