package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePageRender(t *testing.T) {
	t.Parallel()

	browser := GetBrowser()
	defer browser.Close()

	page := NewHomePage(browser)
	assert.NotNil(t, page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage(WikiUriBase + "home.png")
}

func TestClickCreateAGame(t *testing.T) {
	t.Parallel()
	browser := GetBrowser()
	defer browser.Close()

	page := NewHomePage(browser)
	createAGameButton := page.GetCreateGameButton()
	assert.NotNil(t, createAGameButton, "Should find a create a game button")

	createAGameButton.MustClick()

	page.Page.MustWaitStable()
	assert.Equal(t, GetBasePage()+"game/create", page.Page.MustInfo().URL, "Should go to the create page on click")
}
