package main

import (
	"log"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

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

type WithServicesSuite struct {
	suite.Suite
	backendProcess  *exec.Cmd
	frontendProcess *exec.Cmd
}

func (s *WithServicesSuite) startFrontend() {
	log.Println("Starting frontend")
	s.frontendProcess = exec.Command("./e2e-start.sh")
	err := s.frontendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start frontend: %s", err)
	}
	log.Println("Started frontend")
}

func (s *WithServicesSuite) startBackend() {
	log.Println("Starting backend")
	s.backendProcess = exec.Command("../backend/backend")
	err := s.backendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start backend: %s", err)
	}
	log.Println("Started backend")
}

func (s *WithServicesSuite) SetupSuite() {
	log.Printf("Starting services")
	s.startBackend()
	s.startFrontend()

	log.Printf("Waiting for services to become ready")
	waitForBackend()
	waitForFrontend()
}

func (s *WithServicesSuite) TearDownSuite() {
	log.Print("Shutting down test services")
	if err := s.backendProcess.Process.Kill(); err != nil {
		log.Printf("Cannot kill backend: %s", err)
	}
	if err := s.frontendProcess.Process.Kill(); err != nil {
		log.Printf("Cannot kill frontend: %s", err)
	}
}

func TestWithServicesSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WithServicesSuite))
}

func GetBasePage() string {
	return frontendUrl
}
