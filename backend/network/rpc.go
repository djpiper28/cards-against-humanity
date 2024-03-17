package network

import (
	"encoding/json"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
)

type RpcMessageType int

// The type of the message
const (
  // Tx the initial game state when the user joins
	MsgOnJoin = 0
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
