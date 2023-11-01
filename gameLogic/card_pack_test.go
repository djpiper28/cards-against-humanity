package gameLogic_test

import (
	"fmt"
	"os"
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

func TestLoadCards(t *testing.T) {
	os.Chdir("..")
	err := gameLogic.LoadPacks()
	if err != nil {
		t.Log("Error whilst loading cards", err)
		t.FailNow()
	}

	if len(gameLogic.AllBlackCards) == 0 {
		t.Log("There are no black cards")
		t.FailNow()
	}

	if len(gameLogic.AllWhiteCards) == 0 {
		t.Log("There are no white ards")
		t.FailNow()
	}

	if len(gameLogic.AllPacks) == 0 {
		t.Log("There are no packs")
		t.FailNow()
	}
}
