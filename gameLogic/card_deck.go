package gameLogic

import (
	"errors"
	"math/rand"
	"sync"
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

	WhiteCardsCopy := make([]*WhiteCard, len(WhiteCards))
	copy(WhiteCardsCopy, WhiteCards)

	BlackCardsCopy := make([]*BlackCard, len(BlackCards))
	copy(BlackCardsCopy, BlackCards)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		rand.Shuffle(len(WhiteCardsCopy), func(i, j int) {
			WhiteCardsCopy[i], WhiteCardsCopy[j] = WhiteCardsCopy[j], WhiteCardsCopy[i]
		})
		wg.Done()
	}()

	rand.Shuffle(len(WhiteCardsCopy), func(i, j int) {
		BlackCardsCopy[i], BlackCardsCopy[j] = BlackCardsCopy[j], BlackCardsCopy[i]
	})
	wg.Wait()
	return &CardDeck{WhiteCards: WhiteCardsCopy, BlackCards: BlackCardsCopy}, nil
}

func (cd *CardDeck) GetNewWhiteCards(cardsToAdd uint) ([]*WhiteCard, error) {
	if len(cd.WhiteCards) < int(cardsToAdd) {
		return nil, errors.New("There are not enough cards to give to the player")
	}

	cards := cd.WhiteCards[0:cardsToAdd]
	cd.WhiteCards = cd.WhiteCards[cardsToAdd:]
	return cards, nil
}
