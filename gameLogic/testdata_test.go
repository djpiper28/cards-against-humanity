package gameLogic_test

import (
	"fmt"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/google/uuid"
)

const testCardsLength = 100

func GetTestWhiteCards() []*gameLogic.WhiteCard {
	cards := make([]*gameLogic.WhiteCard, testCardsLength)

	for i := 0; i < testCardsLength; i++ {
		cards = append(cards, gameLogic.NewWhiteCard(uuid.New(), fmt.Sprintf("Test white card #%d", i)))
	}

	return cards
}

func GetTestBlackCards() []*gameLogic.BlackCard {
	cards := make([]*gameLogic.BlackCard, testCardsLength)

	for i := 0; i < testCardsLength; i++ {
		cards = append(cards, gameLogic.NewBlackCard(uuid.New(), fmt.Sprintf("Test black card #%d", i), uint(i%5)))
	}

	return cards
}
