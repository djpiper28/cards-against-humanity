package main

import (
	"log"

	"github.com/go-rod/rod"
)

type HomePage struct {
	Page *rod.Page
}

func NewHomePage(b *rod.Browser) HomePage {
	url := GetBasePage()
	log.Printf("Home page: %s", url)
	return HomePage{Page: b.MustPage(url).MustWaitStable()}
}

func (hp *HomePage) GetCreateGameButton() *rod.Element {
	return hp.Page.MustElementR("button", "/create a game/i")
}
