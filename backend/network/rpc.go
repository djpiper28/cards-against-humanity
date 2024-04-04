package network

import (
	"encoding/json"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/google/uuid"
)

type RpcMessageType int

// The type of the message
const (
	// Tx the initial game state when the user joins
	MsgOnJoin = iota + 1
	// Tx the player's name and id when they join
	MsgOnPlayerJoin
	// Tx when the player is created in a game
	MsgOnPlayerCreate
	// Tx the player's id and reason when they disconnect
	MsgOnPlayerDisconnect
)

type RpcMessageBody struct {
	Type RpcMessageType `json:"type"`
	Data any            `json:"data"`
}

type RpcMessage interface {
	Type() RpcMessageType
}

func EncodeRpcMessage(msg RpcMessage) ([]byte, error) {
	body := RpcMessageBody{Type: msg.Type(), Data: msg}
	ret, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type RpcOnJoinMsg struct {
	State gameLogic.GameStateInfo `json:"state"`
}

func (msg RpcOnJoinMsg) Type() RpcMessageType {
	return MsgOnJoin
}

type RpcOnPlayerJoinMsg struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}

func (msg RpcOnPlayerJoinMsg) Type() RpcMessageType {
	return MsgOnPlayerJoin
}

type RpcOnPlayerDisconnectMsg struct {
	Id     uuid.UUID `json:"id"`
	Reason string    `json:"reason"`
}

func (msg RpcOnPlayerDisconnectMsg) Type() RpcMessageType {
	return MsgOnPlayerDisconnect
}

type RpcOnPlayerCreateMsg struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (msg RpcOnPlayerCreateMsg) Type() RpcMessageType {
	return MsgOnPlayerCreate
}
