package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// You must have the dev porxy set up
const frontendUrl = "http://localhost:8000/"

func testFrontend() error {
	_, err := http.Get(frontendUrl)
	return err
}

const maxTries = 10

func waitForFrontend() {
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

func waitForBackend() {
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

type WithServicesSuite struct {
	suite.Suite
	appProcess *exec.Cmd
}

func (s *WithServicesSuite) start() {
	log.Println("Starting frontend")
	s.appProcess = exec.Command("docker-compose", "up", "--build", "--detach")
	s.appProcess.Stdout = os.Stdout
	s.appProcess.Stderr = os.Stderr
	err := s.appProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start frontend: %s", err)
	}
	log.Println("Started frontend")
}

func (s *WithServicesSuite) SetupSuite() {
	log.Printf("Starting services")
	s.start()

	log.Printf("Waiting for services to become ready")
	waitForBackend()
	waitForFrontend()
}

func TestWithServicesSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WithServicesSuite))
}
