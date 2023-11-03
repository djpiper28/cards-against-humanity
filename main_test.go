package main

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const baseUrl = "http://localhost:8080"

func TestServerStart(t *testing.T) {
	go Start()

	t.Log("Sleeping whils the server starts")
	time.Sleep(time.Second)
	resp, err := http.Get(baseUrl + "/healthcheck")
	assert.Nil(t, err, "There should not be an error on the started server", err)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Healthcheck should work")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.Equal(t, `{"healthy":true}`, string(body), "Should return healthy")
}

func TestGetMetrics(t *testing.T) {
	resp, err := http.Get(baseUrl + "/metrics")
	assert.Nil(t, err, "There should not be an error getting the metrics")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.NotEmpty(t, body, "Body should not be empty")
}

func TestGetGamesNotFull(t *testing.T) {
	resp, err := http.Get(baseUrl + "/games/notFull")
	assert.Nil(t, err, "There should not be an error getting the games")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should be OK")

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err, "Should be able to read the body")
	assert.Equal(t, string(body), "[]", "Should be an empty array")
}
