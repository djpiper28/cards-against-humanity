package main

import (
	"log"
	"time"

	"github.com/go-rod/rod"
)

type CreateGamePage struct {
	Page *rod.Page
}

func NewCreateGamePage(b *rod.Browser) CreateGamePage {
	url := GetBasePage() + "create"
	log.Printf("Create Game Page: %s", url)
	createGamePage := CreateGamePage{Page: b.MustPage(url)}

	time.Sleep(Timeout)
	createGamePage.Page.MustWaitStable()
	return createGamePage
}

const (
	createGameInputPlayerName     = "/Player Name/i"
	createGameInputGamePassword   = "/Game Password/i"
	createGameInputMaxPlayers     = "/Max Players/i"
	createGameInputPointsToPlayTo = "/Points To Play To/i"
	createGameInputMaxGameRounds  = "/Max Game Rounds/i"
)

func (c *CreateGamePage) PlayerName() *rod.Element {
	return GetInputByLabel(c.Page, createGameInputPlayerName)
}

func (c *CreateGamePage) GamePasssowrd() *rod.Element {
	return GetInputByLabel(c.Page, createGameInputGamePassword)
}

func (c *CreateGamePage) MaxPlayers() *rod.Element {
	return GetInputByLabel(c.Page, createGameInputMaxPlayers)
}

func (c *CreateGamePage) PointsToPlayTo() *rod.Element {
	return GetInputByLabel(c.Page, createGameInputPointsToPlayTo)
}

func (c *CreateGamePage) MaxGameRounds() *rod.Element {
	return GetInputByLabel(c.Page, createGameInputMaxGameRounds)
}

func (c *CreateGamePage) ClearInputs() {
	c.PlayerName().Input("")
	c.GamePasssowrd().Input("")
	c.MaxPlayers().Input("")
	c.PointsToPlayTo().Input("")
	c.PointsToPlayTo().Input("")
}

// Inserts the default valid settings into the form
func (c *CreateGamePage) InsertDefaultValidSettings() {
	c.ClearInputs()

	c.PlayerName().Input("Steve")
	c.GamePasssowrd().Input("poop")

	GetInputByLabel(c.Page, "/.*CAH:? Base Set.*/").Timeout(Timeout).MustClick()
}

func (c *CreateGamePage) CreateGame() {
	c.Page.Timeout(Timeout).MustElementR("button", "/Create Game/i").MustClick()
}
