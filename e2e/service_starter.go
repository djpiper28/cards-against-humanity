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

	err := process.Run()
	if err != nil {
		log.Fatalln("Cannot start backend")
	}
	started = true
}
