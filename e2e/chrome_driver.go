package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func GetBrowser() *rod.Browser {
	for !isAlive() {
		log.Printf("Sleeping until services are alive...")
		time.Sleep(time.Second / 5)
	}

	var browser *rod.Browser
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		launcher := launcher.New().
			Headless(false).
			NoSandbox(true)

		browser = rod.New().ControlURL(launcher.MustLaunch())
		browser.Trace(true)
		browser.SlowMotion(time.Second / 3)
	} else {
		browser = rod.New()
	}

	return browser.MustConnect()
}

const Timeout = time.Millisecond * 200
const WikiUriBase = "../wiki/assets/"

var errorPngLock sync.Mutex

func screenshotError(p *rod.Page) {
	errorPngLock.Lock()
	defer errorPngLock.Unlock()

	log.Printf("Saving screenshot to error.png")
	p.MustScreenshotFullPage("error.png")
}

func GetInputByLabel(p *rod.Page, label string) *rod.Element {
	defer func() {
		if recover() != nil {
			log.Printf("Error getting input by label '%s', screenshot was saved", label)
			screenshotError(p)
		}
	}()
	return p.Timeout(Timeout).
		MustElementR("label", label).
		MustElement("input")
}

func cssSelectorForId(id string) string {
	if len(id) > 0 {
		if id[0] >= '0' && id[0] <= '9' {
			log.Print("Illegal Id, this causes the CSS selector to be invalid as it starts with a digit")
		}
	}
	return "#" + id
	// return fmt.Sprintf("[id='%s']", id)
}

func GetById(p *rod.Page, id string) *rod.Element {
	defer func() {
		if recover() != nil {
			log.Printf("Error getting element by id '%s', screenshot was saved", id)
			screenshotError(p)
		}
	}()
	return p.Timeout(Timeout).MustElement(cssSelectorForId(id))
}

func GetAllById(p *rod.Page, id string) []*rod.Element {
	defer func() {
		if recover() != nil {
			log.Printf("Error getting all elements by id '%s', screenshot was saved", id)
			screenshotError(p)
		}
	}()
	return p.Timeout(Timeout).MustElements(cssSelectorForId(id))
}

const frontendUrlProxy = "http://localhost:8000/"

func GetBasePage() string {
	return frontendUrlProxy
}

func GetDomain() string {
	url, err := url.Parse(frontendUrlProxy)
	if err != nil {
		log.Fatalf("Cannot get domain for base url: %s", err)
	}
	return url.Host
}

func isAlive() bool {
	_, err := http.Get(frontendUrlProxy)
	if err != nil {
		log.Printf("Cannot talk to the frontend: %s", err)
		return false
	}

	_, err = http.Get(backendUrl + "/res/packs")
	if err != nil {
		log.Printf("Cannot talk to the backend: %s", err)
		return false
	}

	return true
}
