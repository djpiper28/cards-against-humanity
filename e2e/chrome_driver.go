package main

import (
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
)

func GetBrowser() *rod.Browser {
	headless := os.Getenv("HEADLESS") != "false"
	log.Printf("Starting new browser, headless: %v", headless)
	l := launcher.New().
		Leakless(true).
		Headless(headless)

	url := l.MustLaunch()
	b := rod.New().
		MustConnect().
		DefaultDevice(devices.LaptopWithHiDPIScreen).
		ControlURL(url)

	if !headless {
		b.Timeout(time.Second * 5).
			Trace(true)
	}
	return b
}

func GetInputByLabel(p *rod.Page, label string) *rod.Element {
	return p.MustElementR("label", label).MustElement("input")
}
