package main

import (
	"log"

	"github.com/go-rod/rod"
)

type JoinGamePage struct {
	Page     *rod.Page
	GameId   string
	PlayerId string
}

const PlayerIdCookie = "playerId"

func GetJoinGameUrl() string {
	return GetBasePage() + "join"
}

func NewJoinGamePage(b *rod.Browser, gameId string) CreateGamePage {
	url := GetJoinGameUrl() + "?gameId=" + gameId
	log.Printf("Join Game page: %s, gameId: %s",
		url,
		gameId)
	return CreateGamePage{Page: b.MustPage(url).MustWaitStable()}
}