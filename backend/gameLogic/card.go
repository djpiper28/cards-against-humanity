package gameLogic

type WhiteCard struct {
	Id       int    `json:"id"`
	BodyText string `json:"bodyText"`
}

func NewWhiteCard(Id int, BodyText string) *WhiteCard {
	return &WhiteCard{Id: Id, BodyText: BodyText}
}

type BlackCard struct {
	Id          int    `json:"id"`
	BodyText    string `json:"bodyText"`
	CardsToPlay uint   `json:"cardsToPlay"`
}

func NewBlackCard(Id int, BodyText string, CardsToPlay uint) *BlackCard {
	return &BlackCard{Id: Id, BodyText: BodyText, CardsToPlay: CardsToPlay}
}
