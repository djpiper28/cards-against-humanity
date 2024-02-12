package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameRepoInitialised(t *testing.T) {
	assert.NotNil(t, GameRepo, "Game Repo should not be nil")
	assert.NotNil(t, GameRepo.GamesByAge)
	assert.NotNil(t, GameRepo.GameAgeMap)
	assert.NotNil(t, GameRepo.GameMap)
}
