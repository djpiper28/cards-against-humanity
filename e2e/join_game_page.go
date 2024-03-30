package main

import (
	"fmt"
	"log"

	"github.com/go-rod/rod"
)

type JoinGamePage struct {
	Page *rod.Page
}

const PlayerIdCookie = "playerId"

func GetJoinGameUrl() string {
	return GetBasePage() + "game/join"
}

func NewJoinGamePage(b *rod.Browser, gameId string) CreateGamePage {
	url := GetJoinGameUrl() + "?gameId=" + gameId
	log.Printf("Join Game page: %s", url)
	return CreateGamePage{Page: b.MustPage(url).MustWaitStable()}
}

func (page *JoinGamePage) InLobby() bool {
	page.Page.Timeout(Timeout).MustElementR("h1", fmt.Sprintf("/.*'s game/i"))
	return true
}
