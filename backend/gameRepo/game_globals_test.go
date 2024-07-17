package gameRepo_test

import (
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/stretchr/testify/assert"
)

func TestGameRepoInitialised(t *testing.T) {
	assert.NotNil(t, gameRepo.Repo, "Game Repo should not be nil")
	assert.NotNil(t, gameRepo.Repo.GameMap)
}
