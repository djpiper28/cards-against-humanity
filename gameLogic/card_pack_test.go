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

func BenchmarkAccumalateCardPakcs(b *testing.B) {
	packs := GetTestPacks()
	for i := 0; i < b.N; i++ {
		gameLogic.AccumalateCardPacks(packs)
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

func TestBlackCardLookup(t *testing.T) {
	card, err := gameLogic.GetBlackCard(0)
	if err != nil {
		t.Log("Should be able to get black card 0", err)
		t.FailNow()
	}

	if card == nil {
		t.Log("Card was nil")
		t.FailNow()
	}
}

func TestBlackCardLookupNegative(t *testing.T) {
	_, err := gameLogic.GetBlackCard(-1)
	if err == nil {
		t.Log("Should not be able to find a card with ID -1")
		t.FailNow()
	}
}

func TestBlackCardLookupTooHigh(t *testing.T) {
	_, err := gameLogic.GetBlackCard(len(gameLogic.AllBlackCards))
	if err == nil {
		t.Log("Should not be able tof ind a card with ID = len")
		t.FailNow()
	}
}

func TestWhiteCardLookup(t *testing.T) {
	card, err := gameLogic.GetWhiteCard(0)
	if err != nil {
		t.Log("Should be able to get white card 0", err)
		t.FailNow()
	}

	if card == nil {
		t.Log("Card was nil")
		t.FailNow()
	}
}

func TestWhiteCardLookupNegative(t *testing.T) {
	_, err := gameLogic.GetWhiteCard(-1)
	if err == nil {
		t.Log("Should not be able to find a card with ID -1")
		t.FailNow()
	}
}

func TestWhiteCardLookupTooHigh(t *testing.T) {
	_, err := gameLogic.GetWhiteCard(len(gameLogic.AllWhiteCards))
	if err == nil {
		t.Log("Should not be able tof ind a card with ID = len")
		t.FailNow()
	}
}
