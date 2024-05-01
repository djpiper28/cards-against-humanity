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
  assert.True(s.T(), adminLobbyPage.IsAdmin())
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

func (s *WithServicesSuite) TestChangingSettingsSyncsBetweenClients() {
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

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(s.T(), playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{Page: playerPage.Page}

	assert.True(s.T(), playerLobbyPage.InLobby())

	adminLobbyPage.AdminGamePasssowrd().MustInput("Password 123")
	assert.Equal(s.T(), adminLobbyPage.AdminGamePasssowrd().MustText(), "poopPassword 123")

	assert.True(s.T(), adminLobbyPage.Saved())

	found := false
	for _, cookie := range browser.MustGetCookies() {
		if cookie.Name == "password" {
			assert.Equal(s.T(), "poopPassword 123", cookie.Value)
			found = true
			break
		}
	}
	assert.True(s.T(), found)

	found = false
	for _, cookie := range browser2.MustGetCookies() {
		if cookie.Name == "password" {
			assert.Equal(s.T(), "poopPassword 123", cookie.Value)
			found = true
			break
		}
	}
	assert.True(s.T(), found)

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

	playerLobbyPage.Page.MustScreenshotFullPage("../wiki/assets/player_lobby_page.png")
}

func (s *WithServicesSuite) TestPlayerDisconnectReConnect() {
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

	assert.True(s.T(),
		adminLobbyPage.PlayerConnected(adminLobbyPage.PlayerId()))

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(s.T(), playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{Page: playerPage.Page}
	assert.True(s.T(), playerLobbyPage.InLobby())

	playerId := playerLobbyPage.PlayerId()

	time.Sleep(Timeout)
	assert.True(s.T(), adminLobbyPage.PlayerConnected(playerId))

	// Disconnect and make sure the UI updates
	playerPage.Page.MustNavigate("https://google.com").MustActivate()

	time.Sleep(Timeout)
	assert.False(s.T(), adminLobbyPage.PlayerConnected(playerId))

	adminLobbyPage.Page.MustScreenshot("../wiki/assets/player_disconnected.png")

	// Reconnect and make sure the UI updates
	playerPage.Page.MustNavigateBack().MustActivate()
	time.Sleep(Timeout)
	assert.True(s.T(), playerLobbyPage.InLobby())

	time.Sleep(Timeout)
	assert.True(s.T(), adminLobbyPage.PlayerConnected(playerId))
}

func (s *WithServicesSuite) TestPlayerLeavesGame() {
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

	// Connect with another client then assert that the settings remain equal
	browser2 := GetBrowser()
	defer browser2.Close()
	playerPage := NewPlayerGamePage(browser2, adminLobbyPage)

	assert.True(s.T(), playerPage.InPlayerJoinPage())
	playerPage.Password(DefaultPassword)
	playerPage.PlayerName("Geoff")

	playerPage.Join()

	playerLobbyPage := JoinGamePage{Page: playerPage.Page}
	assert.True(s.T(), playerLobbyPage.InLobby())

  playerId := playerLobbyPage.PlayerId()

	playerLobbyPage.LeaveGame()
  assert.NotContains(s.T(), adminLobbyPage.PlayersInGame(), playerId)
  assert.True(s.T(), adminLobbyPage.IsAdmin())
}

// Test that the owner leaving a game transfers ownership to another player
