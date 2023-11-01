package gameLogic_test

import (
	"fmt"
	"testing"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
)

func TestNewCardDeckNilWhiteCards(t *testing.T) {
	blackCards := make([]*gameLogic.BlackCard, 2)

	_, err := gameLogic.NewCardDeck(nil, blackCards)
	if err == nil {
		t.Log("Should not be able to make a deck with no white cards")
		t.FailNow()
	}
}

func TestNewCardDeckNilBlackCards(t *testing.T) {
	whiteCards := make([]*gameLogic.WhiteCard, 2)

	_, err := gameLogic.NewCardDeck(whiteCards, nil)
	if err == nil {
		t.Log("Should not be able to make a deck with no black cards")
		t.FailNow()
	}
}

func TestNewCardDeck(t *testing.T) {
	whiteCards := GetTestWhiteCards()
	blackCards := GetTestBlackCards()

	cardDeck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err != nil {
		t.Log(fmt.Sprintf("Error should be nil %s", err))
		t.FailNow()
	}

	if cardDeck == nil {
		t.Log("Card deck is nil")
		t.FailNow()
	}

	if len(cardDeck.WhiteCards) != len(whiteCards) {
		t.Log("White cards not set")
		t.FailNow()
	}

	if len(cardDeck.BlackCards) != len(blackCards) {
		t.Log("Black cards not set")
		t.FailNow()
	}
}

func TestGetNewWhiteCards(t *testing.T) {
	whiteCards := GetTestWhiteCards()
	testCardsLength := len(whiteCards)
	blackCards := GetTestBlackCards()

	cardDeck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err != nil {
		t.Log(fmt.Sprintf("Error should be nil %s", err))
		t.FailNow()
	}

	if cardDeck == nil {
		t.Log("Card deck is nil")
		t.FailNow()
	}

	var cardsToGet uint = 5
	cards, err := cardDeck.GetNewWhiteCards(uint(cardsToGet))
	if err != nil {
		t.Log(fmt.Sprintf("Error should be nil %s", err))
		t.FailNow()
	}

	if len(cards) != int(cardsToGet) {
		t.Log(fmt.Sprintf("Should have got %d cards, got %d", cardsToGet, len(cards)))
		t.FailNow()
	}

	expected := testCardsLength - int(cardsToGet)
	if len(cardDeck.WhiteCards) != expected {
		t.Log(fmt.Sprintf("Not enough cards where removed %d should be %d", len(cardDeck.WhiteCards), expected))
		t.FailNow()
	}
}

func TestCardDeckGetNewBlackCard(t *testing.T) {
  whiteCards := GetTestWhiteCards()
  blackCards := GetTestBlackCards()

  cardDeck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
  if err != nil {
    t.Log("Should have got a new black card")
    t.FailNow()
  }

  if len(cardDeck.BlackCards) != testCardsLength - 1 {
    t.Log("A black card was not removed")
    t.FailNow()
  }
}

func TestCardDeckGetNewBlackCardNoneLeft(t *testing.T) {
  whiteCards := GetTestWhiteCards()
  blackCards := []*gameLogic.BlackCard{}

  cardDeck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
  if err != nil {
    t.Log("An error occurred making the card deck", err)
    t.FailNow()
  }

  _, err = cardDeck.GetNewBlackCard()
  if err == nil {
    t.Log("Should not have been able to get a new black card")
    t.FailNow()
  }
}

func TestCardDeckAccumalate(t *testing.T) {
	deckCount := 10
	decks := make([]*gameLogic.CardDeck, deckCount)
	for i := 0; i < deckCount; i++ {
		whiteCards := GetTestWhiteCards()
		blackCards := GetTestBlackCards()
		cardDeck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
		if err != nil {
			t.Log("Cannot create the test card deck")
			t.FailNow()
		}

		decks[i] = cardDeck
	}

	accDeck, err := gameLogic.AccumalateDecks(decks)
	if err != nil {
		t.Log("Error accumalting cards should be nil", err)
		t.FailNow()
	}

	expectedWhiteCards := testCardsLength * deckCount
	if len(accDeck.WhiteCards) != expectedWhiteCards {
		t.Log(fmt.Sprintf("Expeceted %d cards. found %d", expectedWhiteCards, len(accDeck.WhiteCards)))
		t.FailNow()
	}

	expectedBlackCards := testCardsLength * deckCount
	if len(accDeck.BlackCards) != expectedBlackCards {
		t.Log(fmt.Sprintf("Expeceted %d cards. found %d", expectedBlackCards, len(accDeck.BlackCards)))
		t.FailNow()
	}
}

func TestCardAccumlateNoWhiteCards(t *testing.T) {
	whiteCards := make([]*gameLogic.WhiteCard, 0)
	blackCards := make([]*gameLogic.BlackCard, 20)
	deck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err != nil {
		t.Log("Should be able to make a deck with no white cards")
		t.FailNow()
	}

	decks := []*gameLogic.CardDeck{deck}
	_, err = gameLogic.AccumalateDecks(decks)
	if err == nil {
		t.Log("Should error when there are no white cards in resultant deck")
		t.FailNow()
	}
}

func TestCardAccumlateNoBlackCards(t *testing.T) {
	whiteCards := make([]*gameLogic.WhiteCard, 20)
	blackCards := make([]*gameLogic.BlackCard, 0)
	deck, err := gameLogic.NewCardDeck(whiteCards, blackCards)
	if err != nil {
		t.Log("Should be able to make a deck with no black cards")
		t.FailNow()
	}

	decks := []*gameLogic.CardDeck{deck}
	_, err = gameLogic.AccumalateDecks(decks)
	if err == nil {
		t.Log("Should error when there are no black cards in resultant deck")
		t.FailNow()
	}
}
