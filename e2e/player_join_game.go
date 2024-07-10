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

func (p *PlayerJoinGame) InLobbyPlayer() bool {
	if !strings.Contains(p.Page.MustInfo().URL, GetJoinGameUrl()) {
		log.Printf("Cannot be in lobby - not under %s", GetJoinGameUrl())
		return false
	}
	return GetById(p.Page, "leave-game") != nil
}

func (p *PlayerJoinGame) PlayerName(name string) {
	GetInputByLabel(p.Page, "/Player Name/i").MustInput(name)
}

func (p *PlayerJoinGame) Password(password string) {
	GetInputByLabel(p.Page, "/Password, leave blank if none/i").MustInput(password)
}

func (p *PlayerJoinGame) Join() {
	p.Page.Timeout(Timeout).MustElementR("button", "/Join Game/i").MustClick()
	p.Page.MustWaitStable()
	return
}

func (p *PlayerJoinGame) Disconnect() {
	p.Page.MustNavigate("https://google.com").MustWaitStable()
}

func (p *PlayerJoinGame) ReConnect() {
	p.Page.Timeout(Timeout * 5).MustNavigate(GetJoinGameUrl()).MustWaitStable()
}

type Card struct {
	Id   string
	Text string
}

func (p *PlayerJoinGame) Cards() ([]Card, error) {
	cards := make([]Card, 0)
	for _, node := range p.Page.Timeout(Timeout).MustElement(cssSelectorForId("card-list")).MustDescribe().Children {
		el := p.Page.MustElementFromNode(node)

		cards = append(cards, Card{
			Id:   *el.MustAttribute("id"),
			Text: el.MustElement("p").MustText(),
		})
	}

	if len(cards) == 0 {
		return nil, errors.New("No cards found")
	}
	return cards, nil
}
