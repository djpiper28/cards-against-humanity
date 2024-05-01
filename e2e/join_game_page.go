package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-rod/rod"
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

const (
	joinGameInputPlayerName     = "/Player Name/i"
	joinGameInputGamePassword   = "/Game Password/i"
	joinGameInputMaxPlayers     = "/Max Players/i"
	joinGameInputPointsToPlayTo = "/Points To Play To/i"
	joinGameInputMaxGameRounds  = "/Max Game Rounds/i"
)

func (j *JoinGamePage) AdminPlayerName() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputPlayerName)
}

func (j *JoinGamePage) AdminGamePasssowrd() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputGamePassword)
}

func (j *JoinGamePage) AdminMaxPlayers() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputMaxPlayers)
}

func (j *JoinGamePage) AdminPointsToPlayTo() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputPointsToPlayTo)
}

func (j *JoinGamePage) AdminMaxGameRounds() *rod.Element {
	return GetInputByLabel(j.Page, joinGameInputMaxGameRounds)
}

func (j *JoinGamePage) Saved() bool {
	return j.Page.Timeout(Timeout).MustElementR("p", "/Settings are saved./i") != nil
}

const (
	joinGameViewGamePasswordId    = "game-password"
	joinGameViewMaxPlayersId      = "max-players"
	joinGameViewPlayingToPointsId = "playing-to-points"
	joinGameViewMaxGameRoundsId   = "max-game-rounds"
)

func (j *JoinGamePage) UserGamePassword() *rod.Element {
	return j.Page.Timeout(Timeout).MustElement("#" + joinGameViewGamePasswordId)
}

func (j *JoinGamePage) UserMaxPlayers() *rod.Element {
	return j.Page.Timeout(Timeout).MustElement("#" + joinGameViewMaxPlayersId)
}

func (j *JoinGamePage) UserPlayingToPoints() *rod.Element {
	return j.Page.Timeout(Timeout).MustElement("#" + joinGameViewPlayingToPointsId)
}

func (j *JoinGamePage) UserMaxGameRounds() *rod.Element {
	return j.Page.Timeout(Timeout).MustElement("#" + joinGameViewMaxGameRoundsId)
}

func (j *JoinGamePage) HasCardPack(packId string) bool {
	return j.Page.Timeout(Timeout).MustElement("#"+packId) != nil
}

func (j *JoinGamePage) PlayerId() string {
	for _, cookie := range j.Page.MustCookies() {
		if cookie.Name == PlayerIdCookie {
			return cookie.Value
		}
	}
	log.Println("The player has no id, this is probably bad")
	return ""
}

func (j *JoinGamePage) PlayerConnected(playerId string) bool {
	defer func() {
		if pan := recover(); pan != nil {
			log.Printf("Error checking if player %s is connected, %s was saved", playerId, ErrorScreenshot)
			j.Page.MustScreenshot(ErrorScreenshot)
			panic(pan)
		}
	}()

	domId := fmt.Sprintf("p#%s-player-status", playerId)
	el := j.Page.Timeout(Timeout).MustElement(domId)
	return strings.ToLower(el.MustText()) == "connected"
}
