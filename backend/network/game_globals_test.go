package network_test

import (
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/stretchr/testify/assert"
)

func TestGameRepoInitialised(t *testing.T) {
	assert.NotNil(t, network.GameRepo, "Game Repo should not be nil")
	assert.NotNil(t, network.GameRepo.GamesByAge)
	assert.NotNil(t, network.GameRepo.GameAgeMap)
	assert.NotNil(t, network.GameRepo.GameMap)
}
