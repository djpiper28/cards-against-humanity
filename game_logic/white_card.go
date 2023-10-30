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
