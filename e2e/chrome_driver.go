package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
)

func GetBasePage() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get executable: %s", err)
	}

	return "file://" + filepath.Dir(wd) + "/cahfrontend/dist/"
}

func GetBrowser() *rod.Browser {
	return rod.New().MustConnect()
}
