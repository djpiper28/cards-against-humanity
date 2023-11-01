package gameLogic

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type Player struct {
	Id          uuid.UUID
	Name        string
	Hand        map[uuid.UUID]*WhiteCard
	CurrentPlay []*WhiteCard
	Connected   bool
	Points      int
}

const (
	maxLength = 20
	minLength = 3
)

func NewPlayer(Name string) (*Player, error) {
	if len(Name) > maxLength || len(Name) < minLength {
		return nil, errors.New(fmt.Sprintf("Length of name must be between %d and %d (exclusive exclusive)", minLength, maxLength))
	}

	return &Player{Id: uuid.New(),
		Name:      Name,
		Hand:      make(map[uuid.UUID]*WhiteCard),
		Connected: true}, nil
}

func (p *Player) hasCard(card *WhiteCard) bool {
	_, found := p.Hand[card.Id]
	return found
}

func (p *Player) PlayCard(cards []*WhiteCard) error {
	if cards == nil {
		return errors.New("Cannot play nil cards")
	}

	if p.CurrentPlay != nil {
		return errors.New("Cards have already been played")
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
	if p.hasCard(card) {
		msg := "Cannot add duplicate cards to the hand"
		log.Println(msg)
		return errors.New(msg)
	}

	p.Hand[card.Id] = card
	return nil
}

func (p *Player) FinaliseRound() {
	p.CurrentPlay = nil
}
