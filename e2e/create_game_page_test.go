package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGamePageRender(t *testing.T) {
	t.Parallel()
	page := NewCreatePagePage(GetBrowser())
	assert.NotNil(t, page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage("../wiki/assets/create_game.png")
}

func TestCreateGamePageDefaultInput(t *testing.T) {
	t.Parallel()
	page := NewCreatePagePage(GetBrowser())
	assert.NotNil(t, page, "Page should render and not be nil")
	page.InsertDefaultValidSettings()
	page.Page.MustScreenshotFullPage("../wiki/assets/create_game_default_input.png")
}
