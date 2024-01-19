package main

import (
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

var lock sync.Mutex
var started = false
var backendProcess *exec.Cmd = exec.Command("../backend/backend")
var frontendProcess *exec.Cmd = exec.Command("./e2e-start.sh")

const frontendUrl = "http://localhost:3000/"

func testFrontend() error {
	_, err := http.Get(frontendUrl)
	return err
}

func waitForFrontend() {
	const maxTries = 10
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

const backendUrl = "http://localhost:8080/"

func testBackend() error {
	_, err := http.Get(backendUrl)
	return err
}

func waitForBackend() {
	const maxTries = 10
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

func startFrontend() {
	log.Println("Starting frontend")
	err := frontendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start frontend: %s", err)
	}
	log.Println("Started frontend")
}

func startBackend() {
	log.Println("Starting backend")
	err := backendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start backend: %s", err)
	}
	log.Println("Started backend")
}

func StartService() {
	lock.Lock()
	defer lock.Unlock()
	if started {
		return
	}

	startBackend()
	startFrontend()

	started = true

	waitForBackend()
	waitForFrontend()
}

func GetBasePage() string {
	return frontendUrl
}
