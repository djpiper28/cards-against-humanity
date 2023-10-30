package game_logic_test

import (
	"github.com/djpiper28/cards-against-humanity/game_logic"
	"github.com/google/uuid"
	"testing"
)

func TestPlayerInit(t *testing.T) {
	name := "Dave"
	player := game_logic.NewPlayer(name)

	if player == nil {
		t.Log("Player is nil")
		t.FailNow()
	}

	if player.Name != name {
		t.Log("Player name is not set")
		t.FailNow()
	}

	var nilId uuid.UUID
	if nilId == player.Id {
		t.Log("Player ID is nil")
		t.FailNow()
	}

	if player.Hand == nil {
		t.Log("Player hand is nil")
		t.FailNow()
	}

	if len(player.Hand) != 0 {
		t.Log("Player hand should be empty")
		t.FailNow()
	}

	if player.CurrentPlay != nil {
		t.Log("Player current play isn't nil")
		t.FailNow()
	}

	if !player.Connected {
		t.Log("Player is not connected")
		t.FailNow()
	}
}
