package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBrowser(t *testing.T) {
	t.Parallel()

	browser := GetBrowser()
	defer browser.Close()

	assert.NotNil(t, browser, "Browser must not be nil")
}
