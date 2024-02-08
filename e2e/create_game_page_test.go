package main

import (
	"github.com/stretchr/testify/assert"
)

func (s *WithServicesSuite) TestCreateGamePageRender() {
	page := NewCreatePagePage(GetBrowser())
	assert.NotNil(s.T(), page, "Page should render and not be nil")
	page.Page.MustScreenshotFullPage("../wiki/assets/create_game.png")
}

func (s *WithServicesSuite) TestCreateGamePageDefaultInput() {
	page := NewCreatePagePage(GetBrowser())
	assert.NotNil(s.T(), page, "Page should render and not be nil")
	page.InsertDefaultValidSettings()
	page.Page.MustScreenshotFullPage("../wiki/assets/create_game_default_input.png")
}
