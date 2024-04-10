package main

import (
	"github.com/stretchr/testify/assert"
)

func (s *WithServicesSuite) TestHomePageRender() {
	browser := GetBrowser()
	defer browser.Close()

	page := NewHomePage(browser)
	assert.NotNil(s.T(), page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage("../wiki/assets/home.png")
}

func (s *WithServicesSuite) TestClickCreateAGame() {
	browser := GetBrowser()
	defer browser.Close()

	page := NewHomePage(browser)
	createAGameButton := page.GetCreateGameButton()
	assert.NotNil(s.T(), createAGameButton, "Should find a create a game button")

	createAGameButton.MustClick()

	page.Page.MustWaitStable()
	assert.Equal(s.T(), GetBasePage()+"game/create", page.Page.MustInfo().URL, "Should go to the create page on click")
}
