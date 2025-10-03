package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-rod/rod"
)

// The page for joining a game, when trying to access admin settings in the lobby use this page
type JoinGamePage struct {
	PlayerJoinGame
}

const PlayerIdCookie = "playerId"

func GetJoinGameUrl() string {
	return GetBasePage() + "game/join"
}

func NewJoinGamePage(b *rod.Browser, gameId string) JoinGamePage {
	url := GetJoinGameUrl() + "?gameId=" + gameId
	log.Printf("Join Game page: %s", url)
	return JoinGamePage{PlayerJoinGame{Page: b.MustPage(url).MustWaitStable()}}
}

func (page *JoinGamePage) InLobbyAdmin() bool {
	page.Page.Timeout(Timeout).MustElementR("h1", "/.*'s game/i")
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
	domId := fmt.Sprintf("player-status-%s", playerId)
	el := GetById(j.Page, domId)
	return strings.ToLower(el.MustText()) == "connected"
}

func (j *JoinGamePage) PlayersInGame() []string {
	players := []string{}
	for _, el := range GetAllById(j.Page, "player-list") {
		players = append(players, el.MustElements("p")[0].MustText())
	}
	return players
}

func (j *JoinGamePage) LeaveGame() {
	GetById(j.Page, "leave-game").Timeout(Timeout).MustClick()
}

func (j *JoinGamePage) IsAdmin() bool {
	return GetById(j.Page, "start-game") != nil
}

func (p *JoinGamePage) Start() {
	p.Page.Timeout(Timeout).MustElement(cssSelectorForId("start-game")).MustClick()
}
