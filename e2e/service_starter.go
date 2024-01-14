package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"
)

var started = false
var backendProcess *exec.Cmd = exec.Command("../backend/backend")
var frontendProcess *exec.Cmd = exec.Command("../cahfrontend/e2e-start.sh")

func StartService() {
	if started {
		return
	}

  log.Println("Starting backend")
	err := backendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start backend: %s", err)
	}
	log.Println("Started backend")

  log.Println("Starting frontend")
	err = frontendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start frontend: %s", err)
	}
	log.Println("Started frontend")

	started = true
}


const frontendUrl = "http://localhost:3000/"
func testFrontend() error {
  _, err := http.Get(frontendUrl)
  return err
}

func GetBasePage() string {
  const maxTries = 10
  for tries := 0; tries < maxTries; tries++ {
    err := testFrontend()
    if err != nil {
      log.Printf("Try %d/%d failed: %s", tries + 1, maxTries, err)
    } else {
      break
    }

    time.Sleep(time.Second)
  }

  return frontendUrl
}
