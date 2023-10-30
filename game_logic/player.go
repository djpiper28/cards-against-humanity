package game_logic

import (
	"github.com/google/uuid"
)

type Player struct {
	Id          uuid.UUID
	Name        string
	Hand        []*WhiteCard
	CurrentPlay *WhiteCard
	Connected   bool
}

func NewPlayer(Name string) *Player {
	return &Player{Id: uuid.New(),
		Name:      Name,
		Hand:      make([]*WhiteCard, 0),
		Connected: true}
}
