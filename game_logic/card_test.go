package game_logic_test

import (
	"github.com/djpiper28/cards-against-humanity/game_logic"
	"github.com/google/uuid"
	"testing"
)

func TestWhiteCardNew(t *testing.T) {
	id := uuid.New()
	bodyText := "bodyText"
	card := game_logic.NewWhiteCard(id, bodyText)

	if card == nil {
		t.Log("Card was null")
		t.FailNow()
	}

	if card.BodyText != bodyText {
		t.Log("Body text not set")
		t.FailNow()
	}

	if card.Id != id {
		t.Log("Id not set")
		t.FailNow()
	}
}

func TestBlackCardNew(t *testing.T) {
	id := uuid.New()
	bodyText := "bodyText"
	var cardsToPlay uint = 12

	card := game_logic.NewBlackCard(id, bodyText, cardsToPlay)

	if card == nil {
		t.Log("Card was null")
		t.FailNow()
	}

	if card.BodyText != bodyText {
		t.Log("Body text not set")
		t.FailNow()
	}

	if card.Id != id {
		t.Log("Id not set")
		t.FailNow()
	}

	if card.CardsToPlay != cardsToPlay {
		t.Log("Cards to play not set")
		t.FailNow()
	}
}
