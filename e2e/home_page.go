package main

import (
	"log"

	"github.com/go-rod/rod"
)

func HomePage(b *rod.Browser) *rod.Page {
	url := GetBasePage() + "index.html"
	log.Printf("Home page: %s", url)
	return b.MustPage(url).MustWaitStable()
}
