package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePageRender(t *testing.T) {
	page := NewHomePage(GetBrowser())
	assert.NotNil(t, page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage("../wiki/assets/home.png")
}

func TestClickCreateAGame(t *testing.T) {
	page := NewHomePage(GetBrowser())
	createAGameButton := page.GetCreateGameButton()
	assert.NotNil(t, createAGameButton, "Should find a create a game button")

	createAGameButton.MustClick()

	page.Page.MustWaitStable()
	assert.Equal(t, GetBasePage()+"create", page.Page.MustInfo().URL, "Should go to the create page on click")
}
