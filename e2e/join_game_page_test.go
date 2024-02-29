package main

import (
	"time"

	"github.com/stretchr/testify/assert"
)

const timeout = time.Millisecond * 200

func (s *WithServicesSuite) TestJoinGameRedirectsOnEmptyGameId() {
	t := s.T()
	t.Parallel()

	joinGame := NewJoinGamePage(GetBrowser(), "")
	time.Sleep(timeout)
	assert.Equal(t, GetBasePage(), joinGame.Page.MustInfo().URL)
}

func (s *WithServicesSuite) TestJoinGameRedirectsOnEmptyPlayerId() {
	t := s.T()
	t.Parallel()

	gameId := "testing123"
	joinGame := NewJoinGamePage(GetBrowser(), gameId)
	time.Sleep(timeout)
	assert.Equal(t, GetBasePage()+"player-join?gameId="+gameId, joinGame.Page.MustInfo().URL)
}