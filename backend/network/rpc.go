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

	// Rx when a player plays a card
	MsgPlayCards
	// Tx to announce that a player has played a card
	MsgOnCardPlayed
	// Tx when the judging phase starts
	MsgOnCzarJudgingPhase

	// Rx when a czar selects a card
	MsgCzarSelectCard
	// Tx when the game state returns to white cards being played
	MsgOnWhiteCardPlayPhase
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
	PlayCardsHandler      func(msg RpcPlayCardsMsg) error
	CzarSelectCardHandler func(msg RpcCzarSelectCardMsg) error
}

func DecodeAs[T RpcMessage](data []byte) (T, error) {
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
		command, err := DecodeAs[RpcChangeSettingsMsg](data)
		if err != nil {
			return err
		}

		return handlers.ChangeSettingsHandler(command)
	case MsgPing:
		// The ping has no body so we don't bother to check it
		return handlers.PingHandler()
	case MsgStartGame:
		return handlers.StartGameHandler()
	case MsgPlayCards:
		command, err := DecodeAs[RpcPlayCardsMsg](data)
		if err != nil {
			return err
		}

		return handlers.PlayCardsHandler(command)
	case MsgCzarSelectCard:
		command, err := DecodeAs[RpcCzarSelectCardMsg](data)
		if err != nil {
			return err
		}

		return handlers.CzarSelectCardHandler(command)
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
	YourPlays         []gameLogic.WhiteCard `json:"yourPlays"`
	// Total amuont of players who have played cards, including yourself:
	// if playerCount - 1 /*card czar*/ == TotalPlays then all players have played
	TotalPlays int `json:"totalPlays"`
}

func (msg RpcRoundInformationMsg) Type() RpcMessageType {
	return MsgRoundInformation
}

type RpcPlayCardsMsg struct {
	CardIds []int `json:"cardIds"`
}

func (msg RpcPlayCardsMsg) Type() RpcMessageType {
	return MsgPlayCards
}

type RpcOnCardPlayedMsg struct {
	PlayerId uuid.UUID `json:"playerId"`
}

func (msg RpcOnCardPlayedMsg) Type() RpcMessageType {
	return MsgOnCardPlayed
}

type RpcOnCzarJudgingPhaseMsg struct {
	// An anonymous list of all of the plays that have been made by all players
	AllPlays [][]*gameLogic.WhiteCard `json:"allPlays"`
	// Your new hand (player specific), if you are the card czar, this is just your hand.
	NewHand []*gameLogic.WhiteCard `json:"newHand"`
}

func (msg RpcOnCzarJudgingPhaseMsg) Type() RpcMessageType {
	return MsgOnCzarJudgingPhase
}

type RpcCzarSelectCardMsg struct {
	// Unsorted array of card IDs, i.e: {1, 3, 2}
	Cards []int `json:"cards"`
}

func (msg RpcCzarSelectCardMsg) Type() RpcMessageType {
	return MsgCzarSelectCard
}

type RpcOnWhiteCardPlayPhase struct {
	BlackCard  *gameLogic.BlackCard   `json:"blackCard"`
	YourHand   []*gameLogic.WhiteCard `json:"yourHand"`
	CardCzarId uuid.UUID              `json:"cardCzarId"`
}

func (msg RpcOnWhiteCardPlayPhase) Type() RpcMessageType {
	return MsgOnWhiteCardPlayPhase
}
