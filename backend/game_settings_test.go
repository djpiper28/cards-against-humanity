package main

import (
	"log"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func (s *ServerTestSuite) TestChangeSettings() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	url := WsBaseUrl + "/games/join"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
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

	// Change the settings
	newSettings := gameLogic.DefaultGameSettings()
	newSettings.Password = "password123"
	newSettings.MaxPlayers = 7

	assert.True(t, newSettings.Validate())
	changeSettingsMsg, err := network.EncodeRpcMessage(network.RpcChangeSettingsMsg{Settings: *newSettings})

	assert.NoError(t, err)
	err = conn.WriteMessage(websocket.TextMessage, changeSettingsMsg)

	assert.NoError(t, err)

	// Check that the settings have been changed
	msgType, msg, err = conn.ReadMessage()

	assert.NoError(t, err)
	assert.Equal(t, msgType, websocket.TextMessage)

	_, err = network.DecodeAs[network.RpcChangeSettingsMsg](msg)
	assert.NoError(t, err)
	assert.Equal(t, changeSettingsMsg, msg)
}
