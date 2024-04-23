package gameLogic_test

import (
	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/google/uuid"
	"testing"
)

func TestPlayerInit(t *testing.T) {
	name := "Dave"
	player, err := gameLogic.NewPlayer(name)

	if err != nil {
		t.Log("There should not be an error", err)
		t.FailNow()
	}

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

	if player.Connected {
		t.Log("Player is connected")
		t.FailNow()
	}

	if player.Points != 0 {
		t.Log("Player points should be 0")
		t.FailNow()
	}
}

func TestAddCardToHand(t *testing.T) {
	player, err := gameLogic.NewPlayer("Dave")
	if err != nil {
		t.Log("Should be able to make the player", err)
		t.FailNow()
	}

	id := 123
	err = player.AddCardToHand(gameLogic.NewWhiteCard(id, "testing 123"))
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

func TestPlayHand(t *testing.T) {
	player, err := gameLogic.NewPlayer("Dave")
	if err != nil {
		t.Log("Should be able to make the player", err)
		t.FailNow()
	}

	cards := []*gameLogic.WhiteCard{gameLogic.NewWhiteCard(0, "Testing 123"),
		gameLogic.NewWhiteCard(1, "Testing 234"),
		gameLogic.NewWhiteCard(2, "Testing 345"),
		gameLogic.NewWhiteCard(3, "Testing 456"),
		gameLogic.NewWhiteCard(4, "Testing 567"),
	}

	for _, card := range cards {
		err := player.AddCardToHand(card)
		if err != nil {
			t.Log("Cannot add card to hand")
			t.FailNow()
		}
	}

	play := cards[0:2]
	err = player.PlayCard(play)
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

func TestPlayNilHand(t *testing.T) {
	player, err := gameLogic.NewPlayer("Dave")
	if err != nil {
		t.Log("Should be able to make the player", err)
		t.FailNow()
	}

	err = player.PlayCard(nil)

	if err == nil {
		t.Log("Should not be able to play a nil hand")
		t.FailNow()
	}
}

func TestPlayDuplicateHand(t *testing.T) {
	player, err := gameLogic.NewPlayer("Dave")
	if err != nil {
		t.Log("Should be able to make the player", err)
		t.FailNow()
	}

	cards := []*gameLogic.WhiteCard{gameLogic.NewWhiteCard(0, "Testing 123"),
		gameLogic.NewWhiteCard(1, "Testing 234"),
		gameLogic.NewWhiteCard(2, "Testing 345"),
		gameLogic.NewWhiteCard(3, "Testing 456"),
		gameLogic.NewWhiteCard(4, "Testing 567"),
	}

	for _, card := range cards {
		err := player.AddCardToHand(card)
		if err != nil {
			t.Log("Cannot add card to hand")
			t.FailNow()
		}
	}

	hand := []*gameLogic.WhiteCard{cards[0], cards[1], cards[2], cards[2]}
	err = player.PlayCard(hand)

	if err == nil {
		t.Log("Should not be able to play a nil hand")
		t.FailNow()
	}
}

func TestFinaliseRound(t *testing.T) {
	player, err := gameLogic.NewPlayer("Dave")
	if err != nil {
		t.Log("Should be able to make the player", err)
		t.FailNow()
	}

	player.CurrentPlay = make([]*gameLogic.WhiteCard, 0)
	player.FinaliseRound()

	if player.CurrentPlay != nil {
		t.Log("Current play was not cleared")
		t.FailNow()
	}
}
