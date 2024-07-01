package main

import (
	"fmt"
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

func (s *ServerTestSuite) TestStartGameEnoughPlayers() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	assert.Nil(t, err)
	assert.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	assert.NoError(t, err, "Should be a join message")
	assert.Equal(t, client.GameId, onJoinMsg.State.Id)
	assert.Len(t, onJoinMsg.State.Players, 1)
	assert.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	for i := 1; i < gameLogic.MinPlayers; i++ {
		name := fmt.Sprintf("Player %d", i)
		info, err := client.AddPlayer(name)
		assert.NoError(t, err)

		_, msg, err = client.Read()
		assert.Nil(t, err, "Should be able to read the message")
		assert.True(t, len(msg) > 0, "Message should have a non-zero length")

		rpcMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
		assert.NoError(t, err)
		assert.Equal(t, rpcMsg.Id, info.PlayerId)
		assert.Equal(t, rpcMsg.Name, name)
	}

	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	assert.Nil(t, err)

	err = client.Write(msgBytes)
	assert.NoError(t, err)

	_, msg, err = client.Read()
	assert.Nil(t, err, "Should be able to read the message")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcRoundInformationMsg](msg)
	assert.NoError(t, err)
	assert.Len(t, rpcMsg.YourHand, gameLogic.HandSize)
	assert.NotEmpty(t, rpcMsg.BlackCard)
	assert.Empty(t, rpcMsg.YourPlays)
	assert.Equal(t, 0, rpcMsg.TotalPlays)
	assert.Equal(t, uint(1), rpcMsg.RoundNumber)
}

func (s *ServerTestSuite) TestPlayingCardInGame() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	assert.Nil(t, err)
	assert.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	assert.Nil(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	assert.NoError(t, err, "Should be a join message")
	assert.Equal(t, client.GameId, onJoinMsg.State.Id)
	assert.Len(t, onJoinMsg.State.Players, 1)
	assert.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	for i := 1; i < gameLogic.MinPlayers; i++ {
		name := fmt.Sprintf("Player %d", i)
		info, err := client.AddPlayer(name)
		assert.NoError(t, err)

		_, msg, err = client.Read()
		assert.Nil(t, err, "Should be able to read the message")
		assert.True(t, len(msg) > 0, "Message should have a non-zero length")

		rpcMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
		assert.NoError(t, err)
		assert.Equal(t, rpcMsg.Id, info.PlayerId)
		assert.Equal(t, rpcMsg.Name, name)
	}

	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	assert.Nil(t, err)

	err = client.Write(msgBytes)
	assert.NoError(t, err)

	_, msg, err = client.Read()
	assert.Nil(t, err, "Should be able to read the message")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcRoundInformationMsg](msg)
	assert.NoError(t, err)
	assert.Len(t, rpcMsg.YourHand, gameLogic.HandSize)
	assert.NotEmpty(t, rpcMsg.BlackCard)
	assert.Empty(t, rpcMsg.YourPlays)
	assert.Equal(t, 0, rpcMsg.TotalPlays)
	assert.Equal(t, uint(1), rpcMsg.RoundNumber)

	playCardMsg, err := network.EncodeRpcMessage(network.RpcPlayCardsMsg{CardIds: []int{rpcMsg.YourHand[0].Id}})
	assert.NoError(t, err)

	err = client.Write(playCardMsg)
	assert.NoError(t, err)

	_, msg, err = client.Read()
	assert.NoError(t, err)

	cardPlayedMsg, err := network.DecodeAs[network.RpcOnCardPlayedMsg](msg)
	assert.NoError(t, err)
	assert.Equal(t, client.PlayerId, cardPlayedMsg.PlayerId)
}
