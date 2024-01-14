package main

import (
	"github.com/go-rod/rod"
)

func GetBrowser() *rod.Browser {
	StartService()
	return rod.New().MustConnect()
}
