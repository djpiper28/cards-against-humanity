package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBrowser(t *testing.T) {
	browser := GetBrowser()
	assert.NotNil(t, browser, "Browser must not be nil")
}
