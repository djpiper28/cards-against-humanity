package gameLogic_test

import (
	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"testing"
)

func TestWhiteCardNew(t *testing.T) {
	id := 123
	bodyText := "bodyText"
	card := gameLogic.NewWhiteCard(id, bodyText)

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
	id := 123
	bodyText := "bodyText"
	var cardsToPlay uint = 12

	card := gameLogic.NewBlackCard(id, bodyText, cardsToPlay)

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
