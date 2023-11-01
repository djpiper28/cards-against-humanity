package gameLogic

import (
	"errors"

	"github.com/google/uuid"
)

type CardPack struct {
	Id       uuid.UUID
	Name     string
	CardDeck *CardDeck
}

func AccumalateCardPacks(packs []*CardPack) (*CardDeck, error) {
	if len(packs) == 0 {
		return nil, errors.New("At least one card pack must be selected")
	}

	decks := make([]*CardDeck, len(packs))
	for i, pack := range packs {
		decks[i] = pack.CardDeck
	}

	return AccumlateDecks(decks), nil
}
