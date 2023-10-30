package game_logic

import (
	"github.com/google/uuid"
)

type WhiteCard struct {
	Id       uuid.UUID `json:"id"`
	BodyText string    `json:"bodyText"`
}

func NewWhiteCard(Id uuid.UUID, BodyText string) *WhiteCard {
	return &WhiteCard{Id: Id, BodyText: BodyText}
}

type BlackCard struct {
	Id          uuid.UUID `json:"id"`
	BodyText    string    `json:"bodyText"`
	CardsToPlay uint      `json:"cardsToPlay"`
}

func NewBlackCard(Id uuid.UUID, BodyText string, CardsToPlay uint) *BlackCard {
	return &BlackCard{Id: Id, BodyText: BodyText, CardsToPlay: CardsToPlay}
}
