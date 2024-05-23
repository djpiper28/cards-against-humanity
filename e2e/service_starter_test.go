package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// You must have the dev porxy set up
const frontendUrl = "http://localhost:3000/"
// Proxy for CORS
const frontendUrlProxy = "http://localhost:3255/"
const Timeout = time.Millisecond * 200

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
	s.frontendProcess.Stdout = os.Stdout
	s.frontendProcess.Stderr = os.Stderr
	err := s.frontendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start frontend: %s", err)
	}
	log.Println("Started frontend")
}

func (s *WithServicesSuite) startBackend() {
	log.Println("Starting backend")
	s.backendProcess = exec.Command("../backend/backend")
	s.backendProcess.Stdout = os.Stdout
	s.backendProcess.Stderr = os.Stderr
	s.backendProcess.Env = append(s.backendProcess.Env, "PORT=8080")
	err := s.backendProcess.Start()
	if err != nil {
		log.Fatalf("Cannot start backend: %s", err)
	}
	log.Println("Started backend")

	go func() {
		err := s.backendProcess.Wait()
		if err != nil {
			log.Printf("Backend exited with err: %s", err)
		}
	}()
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
	log.Print("Shutting down test services...")
	if err := s.backendProcess.Process.Kill(); err != nil {
		log.Printf("Cannot kill backend: %s", err)
	}

	if err := s.frontendProcess.Process.Kill(); err != nil {
		log.Printf("Cannot kill frontend: %s", err)
	}
	log.Print("Shutdown has completed")
}

func TestWithServicesSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WithServicesSuite))
}

func GetBasePage() string {
	return frontendUrlProxy
}

func GetDomain() string {
	url, err := url.Parse(frontendUrlProxy)
	if err != nil {
		log.Fatalf("Cannot get domain for base url: %s", err)
	}
	return url.Host
}
