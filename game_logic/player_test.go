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

func TestAddCardToHand(t *testing.T) {
	player := game_logic.NewPlayer("Dave")
	id := uuid.New()
	err := player.AddCardToHand(game_logic.NewWhiteCard(id, "testing 123"))
	if err != nil {
		t.Log("Card adding failed when it should not have", err)
		t.FailNow()
	}

	if len(player.Hand) != 1 {
		t.Log("Card was not added")
		t.FailNow()
	}

	err = player.AddCardToHand(player.Hand[id])
	if err == nil {
		t.Log("Should not be able to add duplicate cards to the hand")
		t.FailNow()
	}
}

func TestAddCard(t *testing.T) {
	player := game_logic.NewPlayer("Dave")
	cards := []*game_logic.WhiteCard{game_logic.NewWhiteCard(uuid.New(), "Testing 123"),
		game_logic.NewWhiteCard(uuid.New(), "Testing 234"),
		game_logic.NewWhiteCard(uuid.New(), "Testing 345"),
		game_logic.NewWhiteCard(uuid.New(), "Testing 456"),
		game_logic.NewWhiteCard(uuid.New(), "Testing 567"),
	}

	for _, card := range cards {
		err := player.AddCardToHand(card)
		if err != nil {
			t.Log("Cannot add card to hand")
			t.FailNow()
		}
	}

	play := cards[0:2]
	err := player.PlayCard(play)
	if err != nil {
		t.Log("Cannot play cards")
		t.FailNow()
	}

	for _, card := range play {
		_, found := player.Hand[card.Id]
		if found {
			t.Log("Card should have been removed as it was played")
			t.FailNow()
		}
	}

	for _, card := range cards[2:] {
		_, found := player.Hand[card.Id]
		if !found {
			t.Log("Card should not have been removed as it wasn't played")
			t.FailNow()
		}
	}

	if len(player.CurrentPlay) != len(play) {
		t.Log("Current play has wrong length")
		t.FailNow()
	}

	for i, currentPlay := range player.CurrentPlay {
		if currentPlay != cards[i] {
			t.Log("Current play is not what is expected")
			t.FailNow()
		}
	}
}
