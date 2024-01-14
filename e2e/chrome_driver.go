package main

import (
	"github.com/go-rod/rod"
)

func GetBrowser() *rod.Browser {
	return rod.New().MustConnect()
}
