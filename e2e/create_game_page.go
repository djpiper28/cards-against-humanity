package main

import (
	"log"

	"github.com/go-rod/rod"
)

type CreateGamePage struct {
	Page *rod.Page
}

func NewCreateGamePage(b *rod.Browser) CreateGamePage {
	url := GetBasePage() + "create"
	log.Printf("Create Game Page: %s", url)
	return CreateGamePage{Page: b.MustPage(url).MustWaitStable()}
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

// Inserts the default valid settings into the form
func (c *CreateGamePage) InsertDefaultValidSettings() {
	c.PlayerName().MustInput("Steve")
	c.GamePasssowrd().MustInput("poop")
	c.MaxPlayers().MustInput("4")
	c.PointsToPlayTo().MustInput("4")
	c.MaxGameRounds().MustInput("20")
}
