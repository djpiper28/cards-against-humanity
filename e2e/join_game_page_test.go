package main

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *WithServicesSuite) TestJoinGameRedirectsOnEmptyGameId() {
	t := s.T()

	joinGame := NewJoinGamePage(GetBrowser(), "")
	time.Sleep(Timeout)
	assert.Equal(t, GetBasePage(), joinGame.Page.MustInfo().URL)
}

func (s *WithServicesSuite) TestJoinGameRedirectsOnEmptyPlayerId() {
	t := s.T()

	gameId := "testing123"
	joinGame := NewJoinGamePage(GetBrowser(), gameId)
	time.Sleep(Timeout)
	assert.Equal(t, GetBasePage()+"game/playerJoin?gameId="+gameId, joinGame.Page.MustInfo().URL)
}
