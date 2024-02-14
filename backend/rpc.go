package main

import (
	"encoding/json"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
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

func EncodeRpcMessage(msg RpcMessage) ([]byte, error) {
	body := RpcMessageBody{Type: msg.Type()}
	ret, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type RpcOnJoinMsg struct {
	Data gameLogic.GameStateInfo `json:"data"`
}

func (msg RpcOnJoinMsg) Type() RpcMessageType {
	return MsgOnJoin
}
