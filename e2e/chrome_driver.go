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
const Timeout = time.Millisecond * 200

func GetInputByLabel(p *rod.Page, label string) *rod.Element {
	defer func() {
		if recover() != nil {
			log.Printf("Error getting input by label '%s', %s was saved", label, ErrorScreenshot)
			p.MustScreenshotFullPage(ErrorScreenshot)
		}
	}()
	return p.Timeout(Timeout).
		MustElementR("label", label).
		MustElement("input")
}

func cssSelectorForId(id string) string {
	if len(id) > 0 {
		if id[0] >= '0' && id[0] <= '9' {
			log.Fatal("Illegal Id, this causes the CSS selector to be invalid as it starts with a digit")
		}
	}
	return "#" + id
	// return fmt.Sprintf("[id='%s']", id)
}

func GetById(p *rod.Page, id string) *rod.Element {
	defer func() {
		if recover() != nil {
			log.Printf("Error getting element by id '%s', %s was saved", id, ErrorScreenshot)
			p.MustScreenshotFullPage(ErrorScreenshot)
		}
	}()
	return p.Timeout(Timeout).MustElement(cssSelectorForId(id))
}

func GetAllById(p *rod.Page, id string) []*rod.Element {
	defer func() {
		if recover() != nil {
			log.Printf("Error getting all elements by id '%s', %s was saved", id, ErrorScreenshot)
			p.MustScreenshotFullPage(ErrorScreenshot)
		}
	}()
	return p.Timeout(Timeout).MustElements(cssSelectorForId(id))
}
