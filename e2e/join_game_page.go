package main

import (
	"fmt"
	"log"

	"github.com/go-rod/rod"
)

const (
	joinGameInputPlayerName     = "/Player Name/i"
	joinGameInputGamePassword   = "/Game Password/i"
	joinGameInputMaxPlayers     = "/Max Players/i"
	joinGameInputPointsToPlayTo = "/Points To Play To/i"
	joinGameInputMaxGameRounds  = "/Max Game Rounds/i"
)

type JoinGamePage struct {
	Page *rod.Page
}

const PlayerIdCookie = "playerId"

func GetJoinGameUrl() string {
	return GetBasePage() + "game/join"
}

func NewJoinGamePage(b *rod.Browser, gameId string) JoinGamePage {
	url := GetJoinGameUrl() + "?gameId=" + gameId
	log.Printf("Join Game page: %s", url)
	return JoinGamePage{Page: b.MustPage(url).MustWaitStable()}
}

func (page *JoinGamePage) InLobby() bool {
	page.Page.Timeout(Timeout).MustElementR("h1", fmt.Sprintf("/.*'s game/i"))
	return true
}

func (j *JoinGamePage) PlayerName() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputPlayerName)
}

func (j *JoinGamePage) GamePasssowrd() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputGamePassword)
}

func (j *JoinGamePage) MaxPlayers() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputMaxPlayers)
}

func (j *JoinGamePage) PointsToPlayTo() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputPointsToPlayTo)
}

func (j *JoinGamePage) MaxGameRounds() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputMaxGameRounds)
}
