package gameLogic

import (
	"errors"
)

type CardDeck struct {
	WhiteCards []*WhiteCard
	BlackCards []*BlackCard
}

const handSize = 7

func NewCardDeck(WhiteCards []*WhiteCard, BlackCards []*BlackCard) (*CardDeck, error) {
	if WhiteCards == nil || len(WhiteCards) == 0 {
		return nil, errors.New("No white cards")
	}

	if BlackCards == nil || len(BlackCards) == 0 {
		return nil, errors.New("No black cards")
	}

	return &CardDeck{WhiteCards: WhiteCards, BlackCards: BlackCards}, nil
}

func (cd *CardDeck) GetNewWhiteCards(cardsToAdd uint) ([]*WhiteCard, error) {
	if len(cd.WhiteCards) < int(cardsToAdd) {
		return nil, errors.New("There are not enough cards to give to the player")
	}

	cards := cd.WhiteCards[0:cardsToAdd]
	return cards, nil
}
