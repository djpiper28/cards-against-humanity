package main

import (
	"encoding/json"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/google/uuid"
)

type RpcMessageType int

// The type of the message
const (
	MsgOnJoin = 0
)

type RpcMessageBody struct {
	Type RpcMessageType `json:"type"`
}

type RpcMessage interface {
	Type() RpcMessageType
}

func EncodeRpcMessagee(msg RpcMessage) ([]byte, error) {
	body := RpcMessageBody{Type: msg.Type()}
	ret, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type RpcOnJoinMsg struct {
	WhiteCardCount int                    `json:"whiteCardCount"`
	BlackCardCount int                    `json:"blackCardCount"`
	PlayerNames    []string               `json:"playerNames"`
	CardPacks      []uuid.UUID            `json:"cardPacks"`
	GameSettings   gameLogic.GameSettings `json:"GameSettings"`

	BlackCard        *gameLogic.BlackCard `json:"blackCard"`
	WhiteCardsPlayed int                  `json:"WhiteCardsPlayed"`

	CurrentCardCzar string              `json:"currentCardCzar"`
	GameOwner       string              `json:"gameOwner"`
	GameState       gameLogic.GameState `json:"gameState"`
	CurrentRound    uint                `json:"currentRound"`
}

func (msg *RpcOnJoinMsg) Type() RpcMessageType {
	return MsgOnJoin
}
