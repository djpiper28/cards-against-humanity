package main

import (
	"time"

	"github.com/go-rod/rod"
)

func GetBrowser() *rod.Browser {
	StartService()
	return rod.New().MustConnect().Timeout(time.Second * 5)
}

func GetInputByLabel(p *rod.Page, label string) *rod.Element {
	return p.MustElementR("label", label).MustElement("input")
}
