package main

import (
	"log"
	"os/exec"
)

var started = false
var process exec.Cmd = *exec.Command("../backend/backend")

func StartService() {
	if started {
		return
	}

	err := process.Start()
	if err != nil {
		log.Fatalf("Cannot start backend: %s", err)
	}
	log.Println("Started backend")
	started = true
}
