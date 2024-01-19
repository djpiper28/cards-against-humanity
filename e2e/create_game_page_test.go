package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGamePageRender(t *testing.T) {
	page := NewCreatePagePage(GetBrowser())
	assert.NotNil(t, page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage("../wiki/assets/create_game.png")
}
