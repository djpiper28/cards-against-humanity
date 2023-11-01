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
	if WhiteCards == nil {
		return nil, errors.New("White cards cannot be nil")
	}

	if BlackCards == nil {
		return nil, errors.New("Black cards cannot be nil")
	}

	WhiteCardsCopy := make([]*WhiteCard, len(WhiteCards))
	copy(WhiteCardsCopy, WhiteCards)

	BlackCardsCopy := make([]*BlackCard, len(BlackCards))
	copy(BlackCardsCopy, BlackCards)

	deck := &CardDeck{WhiteCards: WhiteCardsCopy, BlackCards: BlackCardsCopy}
	deck.Shuffle()
	return deck, nil
}

func (cd *CardDeck) Shuffle() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		rand.Shuffle(len(cd.WhiteCards), func(i, j int) {
			cd.WhiteCards[i], cd.WhiteCards[j] = cd.WhiteCards[j], cd.WhiteCards[i]
		})
	}()

	rand.Shuffle(len(cd.BlackCards), func(i, j int) {
		cd.BlackCards[i], cd.BlackCards[j] = cd.BlackCards[j], cd.BlackCards[i]
	})
	wg.Wait()
}

func (cd *CardDeck) GetNewWhiteCards(cardsToAdd uint) ([]*WhiteCard, error) {
	if len(cd.WhiteCards) < int(cardsToAdd) {
		return nil, errors.New("There are not enough cards to give to the player")
	}

	cards := cd.WhiteCards[0:cardsToAdd]
	cd.WhiteCards = cd.WhiteCards[cardsToAdd:]
	return cards, nil
}

func AccumalateDecks(decks []*CardDeck) (*CardDeck, error) {
	// Count the cards to preallocate them
	whiteCardsCount := 0
	blackCardsCount := 0
	for _, deck := range decks {
		whiteCardsCount += len(deck.WhiteCards)
		blackCardsCount += len(deck.BlackCards)
	}

	returnDeck := &CardDeck{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		returnDeck.WhiteCards = make([]*WhiteCard, 0, whiteCardsCount)
		for _, deck := range decks {
			returnDeck.WhiteCards = append(returnDeck.WhiteCards, deck.WhiteCards...)
		}
	}()

	returnDeck.BlackCards = make([]*BlackCard, 0, blackCardsCount)
	for _, deck := range decks {
		returnDeck.BlackCards = append(returnDeck.BlackCards, deck.BlackCards...)
	}
	wg.Wait()

	if len(returnDeck.WhiteCards) == 0 {
		return nil, errors.New("No white cards")
	}

	if len(returnDeck.BlackCards) == 0 {
		return nil, errors.New("No black cards")
	}

	returnDeck.Shuffle()
	return returnDeck, nil
}
