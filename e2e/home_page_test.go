package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePageRender(t *testing.T) {
	page := HomePage(GetBrowser())
	assert.NotNil(t, page)
}
