package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-rod/rod"
)

type PlayerJoinGame struct {
	Page *rod.Page
}

// Probably should flow via the join page
func GetPlayerJoinGameUrl() string {
	return GetBasePage() + "game/playerJoin"
}

func UpgradeFromJoinPage(p *JoinGamePage) (PlayerJoinGame, error) {
	ret := PlayerJoinGame{Page: p.Page}
	if !strings.Contains(GetPlayerJoinGameUrl(), ret.Page.MustInfo().URL) {
		logMsg := fmt.Sprintf("Not on the correct page. Expected %s, got %s", GetPlayerJoinGameUrl(), ret.Page.MustInfo().URL)
		log.Print(logMsg)
		return PlayerJoinGame{}, errors.New(logMsg)
	}
	return ret, nil
}

func NewPlayerGamePage(b *rod.Browser, gameId string) PlayerJoinGame {
	return PlayerJoinGame{Page: b.MustPage(GetPlayerJoinGameUrl() + "?gameId=" + gameId).MustWaitStable()}
}
