package main

import (
	"log"

	"github.com/go-rod/rod"
)

func GetBrowser() *rod.Browser {
	return rod.New().
		MustConnect()
}

func GetInputByLabel(p *rod.Page, label string) *rod.Element {
	defer func() {
		const fname = "./error.png"
		if recover() != nil {
			log.Printf("Error getting input by label '%s', %s was saved", label, fname)
			p.MustScreenshotFullPage(fname)
		}
	}()
	return p.Timeout(Timeout).MustElementR("label", label).MustElement("input")
}
