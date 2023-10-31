package gameLogic

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"sync"
)

type Player struct {
	Id          uuid.UUID
	Name        string
	Hand        map[uuid.UUID]*WhiteCard
	CurrentPlay []*WhiteCard
	Connected   bool
	lock        sync.Mutex
}

func NewPlayer(Name string) *Player {
	return &Player{Id: uuid.New(),
		Name:      Name,
		Hand:      make(map[uuid.UUID]*WhiteCard),
		Connected: true}
}

func (p *Player) hasCard(card *WhiteCard) bool {
	_, found := p.Hand[card.Id]
	return found
}

func (p *Player) PlayCard(cards []*WhiteCard) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.CurrentPlay != nil {
		return errors.New("Cannot play two cards")
	}

	for _, card := range cards {
		if !p.hasCard(card) {
			return errors.New("Cannot find the card in the hand")
		}
	}

	cardsCopy := make([]*WhiteCard, len(cards))
	copy(cardsCopy, cards)

	p.CurrentPlay = cardsCopy

	for _, card := range cards {
		delete(p.Hand, card.Id)
	}
	return nil
}

func (p *Player) AddCardToHand(card *WhiteCard) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.hasCard(card) {
		msg := "Cannot add duplicate cards to the hand"
		log.Println(msg)
		return errors.New(msg)
	}

	p.Hand[card.Id] = card
	return nil
}
