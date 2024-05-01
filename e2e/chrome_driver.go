package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func GetBrowser() *rod.Browser {
	var browser *rod.Browser
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		launcher := launcher.New().
			Headless(false)
		browser = rod.New().ControlURL(launcher.MustLaunch())
		browser.Trace(true)
		browser.SlowMotion(time.Second / 3)
	} else {
		browser = rod.New()
	}

	return browser.MustConnect()
}

const ErrorScreenshot = "./error.png"

func GetInputByLabel(p *rod.Page, label string) *rod.Element {
	defer func() {
		if recover() != nil {
			log.Printf("Error getting input by label '%s', %s was saved", label, ErrorScreenshot)
			p.MustScreenshotFullPage(ErrorScreenshot)
		}
	}()
	return p.Timeout(Timeout).MustElementR("label", label).MustElement("input")
}
