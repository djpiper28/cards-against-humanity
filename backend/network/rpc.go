package network

import (
	"encoding/json"
	"errors"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
	"github.com/google/uuid"
)

type RpcMessageType int

// The type of the message
const (
	// Tx the initial game state when the user joins
	MsgOnJoin RpcMessageType = iota + 1
	// Tx the player's name and id when they join
	MsgOnPlayerJoin
	// Tx when the player is created in a game
	MsgOnPlayerCreate
	// Tx the player's id and reason when they disconnect
	MsgOnPlayerDisconnect
	// Tx the player's id and reason when they leave the game
	MsgOnPlayerLeave
	// Tx the player's id when a new owner for a game is selected
	MsgNewOwner

	// Tx when a command cannot be processed
	MsgCommandError

	// Rx change the settings of the game
	MsgChangeSettings
	// Rx & Tx for pinging and "ponging" between the server and client
	MsgPing

	// Rx when the owner starts the game
	MsgStartGame

	// Tx the current round info
	MsgRoundInformation
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

type RpcCommandHandlers struct {
	ChangeSettingsHandler func(msg RpcChangeSettingsMsg) error
	PingHandler           func() error
	StartGameHandler      func() error
}

func decodeAs[T any](data []byte) (T, error) {
	type TProxy struct {
		Type RpcMessageType `json:"type"`
		Data T              `json:"data"`
	}

	var proxy TProxy
	err := json.Unmarshal(data, &proxy)
	if err != nil {
		return proxy.Data, err
	}

	return proxy.Data, nil
}

func DecodeRpcMessage(data []byte, handlers RpcCommandHandlers) error {
	var cmd RpcMessageBody
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		return err
	}

	switch cmd.Type {
	case MsgChangeSettings:
		command, err := decodeAs[RpcChangeSettingsMsg](data)
		if err != nil {
			return err
		}

		return handlers.ChangeSettingsHandler(command)
	case MsgPing:
		// The ping has no body so we don't bother to check it
		return handlers.PingHandler()
	default:
		logger.Logger.Error("Unknown command", "type", cmd.Type)
		return errors.New("Unknown command")
	}
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

type RpcCommandErrorMsg struct {
	Reason string `json:"reason"`
}

func (msg RpcCommandErrorMsg) Type() RpcMessageType {
	return MsgCommandError
}

type RpcChangeSettingsMsg struct {
	Settings gameLogic.GameSettings `json:"settings"`
}

func (msg RpcChangeSettingsMsg) Type() RpcMessageType {
	return MsgChangeSettings
}

type RpcOnPlayerLeaveMsg struct {
	Id     uuid.UUID `json:"id"`
	Reason string    `json:"reason"`
}

func (msg RpcOnPlayerLeaveMsg) Type() RpcMessageType {
	return MsgOnPlayerLeave
}

type RpcNewOwnerMsg struct {
	Id uuid.UUID `json:"id"`
}

func (msg RpcNewOwnerMsg) Type() RpcMessageType {
	return MsgNewOwner
}

type RpcPingMsg struct{}

func (msg RpcPingMsg) Type() RpcMessageType {
	return MsgPing
}

type RpcStartGameMsg struct{}

func (msg RpcStartGameMsg) Type() RpcMessageType {
	return MsgStartGame
}

type RpcRoundInformationMsg struct {
	RoundNumber       uint                  `json:"roundNumber"`
	CurrentCardCzarId uuid.UUID             `json:"currentCardCzarId"`
	BlackCard         gameLogic.BlackCard   `json:"blackCard"`
	YourHand          []gameLogic.WhiteCard `json:"yourHand"`
}

func (msg RpcRoundInformationMsg) Type() RpcMessageType {
	return MsgRoundInformation
}
