package gameLogic_test

import (
	"fmt"
	"testing"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
)

func TestAccumalateCardPacks(t *testing.T) {
	packs := GetTestPacks()
	deck, err := gameLogic.AccumalateCardPacks(packs)

	if err != nil {
		t.Log(fmt.Sprintf("Cannot accumalate packs %s", err))
		t.FailNow()
	}

	expectedCards := testCardsLength * testCardsLength
	if len(deck.WhiteCards) != expectedCards {
		t.Log(fmt.Sprintf("Epected %d white cards, found %d", expectedCards, len(deck.WhiteCards)))
		t.FailNow()
	}

	if len(deck.BlackCards) != expectedCards {
		t.Log(fmt.Sprintf("Epected %d black cards, found %d", expectedCards, len(deck.BlackCards)))
		t.FailNow()
	}
}
