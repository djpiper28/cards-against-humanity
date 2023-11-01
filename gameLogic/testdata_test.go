package gameLogic_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/google/uuid"
)

const testCardsLength = 100

func GetTestWhiteCards() []*gameLogic.WhiteCard {
	cards := make([]*gameLogic.WhiteCard, 0, testCardsLength)

	for i := 0; i < testCardsLength; i++ {
		cards = append(cards, gameLogic.NewWhiteCard(uuid.New(), fmt.Sprintf("Test white card #%d", i)))
	}

	return cards
}

func TestGetTestWhiteCards(t *testing.T) {
	if len(GetTestWhiteCards()) != testCardsLength {
		t.Log("The amount of test white cards was wrong")
		t.FailNow()
	}
}

func TestGetTestBlackCards(t *testing.T) {
	if len(GetTestBlackCards()) != testCardsLength {
		t.Log("The amount of test black cards was wrong")
		t.FailNow()
	}
}

func GetTestBlackCards() []*gameLogic.BlackCard {
	cards := make([]*gameLogic.BlackCard, 0, testCardsLength)

	for i := 0; i < testCardsLength; i++ {
		cards = append(cards, gameLogic.NewBlackCard(uuid.New(), fmt.Sprintf("Test black card #%d", i), uint(i%5)))
	}

	return cards
}

func GetTestPacks() []*gameLogic.CardPack {
	packs := make([]*gameLogic.CardPack, testCardsLength)

	for i := 0; i < len(packs); i++ {
		whiteCards := make([]*gameLogic.WhiteCard, 0, testCardsLength)
		for j := 0; j < testCardsLength; j++ {
			whiteCards = append(whiteCards, gameLogic.NewWhiteCard(uuid.New(), fmt.Sprintf("Test white card #%d-%d", i, j)))
		}

		blackCards := make([]*gameLogic.BlackCard, 0, testCardsLength)
		for j := 0; j < testCardsLength; j++ {
			blackCards = append(blackCards, gameLogic.NewBlackCard(uuid.New(), fmt.Sprintf("Test black card #%d-%d", i, j), uint(j%5)))
		}

		deck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
		if err != nil || deck == nil {
			log.Fatal("Cannot create a deck for the test data")
		}

		packs[i] = &gameLogic.CardPack{Id: uuid.New(), Name: fmt.Sprintf("Test Pack #%d", i), CardDeck: deck}
	}

	return packs
}

func TestGetTestPacks(t *testing.T) {
	packs := GetTestPacks()
	if len(packs) != testCardsLength {
		t.Log("The amount of test packs was wrong")
		t.FailNow()
	}

	for _, pack := range packs {
		if pack == nil {
			t.Log("Pack was nil")
			t.FailNow()
		}

		deck := pack.CardDeck
		if deck == nil {
			t.Log("Deck was nil")
			t.FailNow()
		}

		if len(deck.WhiteCards) != testCardsLength {
			t.Log("Test deck has wrong amount of white cards")
			t.FailNow()
		}

		if len(deck.BlackCards) != testCardsLength {
			t.Log("Test deck hs wrong amount of black cards")
			t.FailNow()
		}
	}
}
