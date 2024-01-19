package main

import (
	"log"

	"github.com/go-rod/rod"
)

type CreateGamePage struct {
	Page *rod.Page
}

func NewCreatePagePage(b *rod.Browser) CreateGamePage {
	url := GetBasePage() + "create"
	log.Printf("Create Game Page page: %s", url)
	return CreateGamePage{Page: b.MustPage(url).MustWaitStable()}
}
