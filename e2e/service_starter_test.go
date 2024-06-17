package main

import (
	"log"
	"net/http"
	"testing"
	"time"
)

// You must have the dev porxy set up
const frontendUrl = "http://localhost:8000/"

func testFrontend() error {
	_, err := http.Get(frontendUrl)
	return err
}

const maxTries = 10

func TestFrontend(t *testing.T) {
	t.Parallel()
	for tries := 0; tries < maxTries; tries++ {
		err := testFrontend()
		if err != nil {
			log.Printf("frontend poll %d/%d failed: %s", tries+1, maxTries, err)
		} else {
			break
		}

		time.Sleep(time.Second)
	}
	log.Println("Frontend responded to poll")
}

const backendUrl = "http://localhost:8002/"

func testBackend() error {
	_, err := http.Get(backendUrl)
	return err
}

func TestBackend(t *testing.T) {
	t.Parallel()

	for tries := 0; tries < maxTries; tries++ {
		err := testBackend()
		if err != nil {
			log.Printf("Backend poll %d/%d failed: %s", tries+1, maxTries, err)
		} else {
			break
		}

		time.Sleep(time.Second)
	}
	log.Println("Backend responded to poll")
}
