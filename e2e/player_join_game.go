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

func UpgradeFromJoinPage(p JoinGamePage) (PlayerJoinGame, error) {
	ret := PlayerJoinGame{Page: p.Page}
	currentUrl := ret.Page.MustInfo().URL

	if !strings.Contains(currentUrl, GetPlayerJoinGameUrl()) {
		logMsg := fmt.Sprintf("Not on the correct page. Expected %s, got %s", GetPlayerJoinGameUrl(), currentUrl)
		log.Print(logMsg)
		return PlayerJoinGame{}, errors.New(logMsg)
	}
	return ret, nil
}

func NewPlayerGamePage(b *rod.Browser, adminJoinPage JoinGamePage) PlayerJoinGame {
	ret := PlayerJoinGame{Page: b.MustPage(adminJoinPage.Page.MustInfo().URL).MustWaitStable()}
	return ret
}

func (p *PlayerJoinGame) InPlayerJoinPage() bool {
	return strings.Contains(p.Page.MustInfo().URL, GetPlayerJoinGameUrl())
}

func (p *PlayerJoinGame) PlayerName(name string) {
	GetInputByLabel(p.Page, "Player Name").MustInput(name)
}

func (p *PlayerJoinGame) Password(password string) {
	GetInputByLabel(p.Page, "Password, leave blank if none").MustInput(password)
}

func (p *PlayerJoinGame) Join() {
	p.Page.Timeout(Timeout).MustElementR("button", "/Join Game/i").MustClick()
	p.Page.MustWaitStable()
	return
}

func (p *PlayerJoinGame) Disconnect() {
	p.Page.Timeout(Timeout).MustNavigate(GetBasePage()).MustActivate()
}

func (p *PlayerJoinGame) ReConnect() {
	p.Page.Timeout(Timeout).MustNavigate(GetJoinGameUrl()).MustActivate()
}
