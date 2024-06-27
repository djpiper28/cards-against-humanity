package main

import (
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func (s *ServerTestSuite) TestStartGameTooFewPlayers() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	url := WsBaseUrl + "/games/join"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	conn, _, err := dialer.Dial(url, game.Jar.Headers())
	assert.Nil(t, err, "Should have connected to the ws server successfully")
	defer conn.Close()
	assert.NotNil(t, conn)

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := conn.ReadMessage()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	assert.Nil(t, err)
	assert.Equal(t, game.Ids.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state

	msgType, msg, err = conn.ReadMessage()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	assert.Nil(t, err, "Should be a join message")
	assert.Equal(t, game.Ids.GameId, onJoinMsg.State.Id)
	assert.Len(t, onJoinMsg.State.Players, 1)
	assert.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        game.Ids.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	assert.Nil(t, err)

	conn.WriteMessage(websocket.TextMessage, msgBytes)
	_, msg, err = conn.ReadMessage()
	assert.Nil(t, err, "Should be able to read the message")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcCommandErrorMsg](msg)
	assert.Nil(t, err, "Should be a command error message")
	assert.NotEmpty(t, rpcMsg.Reason)
}
