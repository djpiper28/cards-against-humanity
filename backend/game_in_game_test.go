package main

import (
	"fmt"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
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
	assert.Empty(t, rpcMsg.PreviousWinnerDetails)
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
	assert.Empty(t, rpcMsg.PreviousWinnerDetails)

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

func (s *ServerTestSuite) TestPlayersGetRoundInfoAfterWinnerSelected() {
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

	assert.NoError(t, err, "Should be able to read (the initial game state)")
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

	// Add players
	for i := 1; i < gameLogic.MinPlayers; i++ {
		name := fmt.Sprintf("Player %d", i)
		info, err := client.AddPlayer(name)
		assert.NoError(t, err)

		_, msg, err = client.Read()
		assert.NoError(t, err, "Should be able to read the message")
		assert.True(t, len(msg) > 0, "Message should have a non-zero length")

		rpcMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
		assert.NoError(t, err)
		assert.Equal(t, rpcMsg.Id, info.PlayerId)
		assert.Equal(t, rpcMsg.Name, name)
	}

	// Start game
	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	assert.Nil(t, err)

	err = client.Write(msgBytes)
	assert.NoError(t, err)

	_, msg, err = client.Read()
	assert.NoError(t, err, "Should be able to read the message")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcRoundInformationMsg](msg)
	assert.NoError(t, err)
	assert.Len(t, rpcMsg.YourHand, gameLogic.HandSize)
	assert.NotEmpty(t, rpcMsg.BlackCard)
	assert.Empty(t, rpcMsg.YourPlays)
	assert.Equal(t, 0, rpcMsg.TotalPlays)
	assert.Equal(t, uint(1), rpcMsg.RoundNumber)
	assert.Empty(t, rpcMsg.PreviousWinnerDetails)

	// Play card
	winningPlay := []int{rpcMsg.YourHand[0].Id}
	playCardMsg, err := network.EncodeRpcMessage(network.RpcPlayCardsMsg{CardIds: winningPlay})
	assert.NoError(t, err)

	err = client.Write(playCardMsg)
	assert.NoError(t, err)

	_, msg, err = client.Read()
	assert.NoError(t, err)

	cardPlayedMsg, err := network.DecodeAs[network.RpcOnCardPlayedMsg](msg)
	assert.NoError(t, err)
	assert.Equal(t, client.PlayerId, cardPlayedMsg.PlayerId)

	// Select the winner
	game, err := gameRepo.Repo.GetGame(client.GameId)
	assert.NoError(t, err)

	oldBlackCard := game.CurrentBlackCard
	whiteCard := game.PlayersMap[client.PlayerId].Hand[0]

	info, err := gameRepo.Repo.CzarSelectsCard(client.GameId, game.CurrentCardCzarId, winningPlay)
	assert.NoError(t, err)

	assert.Equal(t, info.WinnerId, client.PlayerId)
	assert.NotEqual(t, oldBlackCard, info.NewBlackCard)
	assert.False(t, info.GameEnded)
  assert.Contains(t, game.Players, info.NewCzarId)

	for pid, hand := range info.Hands {
		assert.Contains(t, game.Players, pid)
		assert.Len(t, hand, 7)

		for _, card := range hand {
			assert.NotEmpty(t, card)
		}
	}

	expectedPreviousWinner := gameLogic.PreviousWinner{
		PlayerId:  client.PlayerId,
		BlackCard: oldBlackCard,
		Whitecards: []*gameLogic.WhiteCard{
			whiteCard,
		},
	}
	assert.Equal(t, expectedPreviousWinner, game.PreviousWinner)
	assert.Equal(t, expectedPreviousWinner, info.PreviousWinner)

  // The player should have had a RpcOnWhiteCardPlayPhase
	assert.NoError(t, err, "Should be able to read (the initial game state)")
	assert.True(t, len(msg) > 0, "Message should have a non-zero length")
	assert.Equal(t, msgType, websocket.TextMessage)

	whiteCardPlay, err := network.DecodeAs[network.RpcOnWhiteCardPlayPhase](msg)
	assert.NoError(t, err, "Should be a white card play message")
	assert.Equal(t, expectedPreviousWinner, whiteCardPlay.Winner)
	assert.Equal(t, client.PlayerId, whiteCardPlay.WinnerId)
  assert.Equal(t, game.CurrentBlackCard, whiteCardPlay.BlackCard)
  assert.Len(t, whiteCardPlay.YourHand, 7)
  assert.Equal(t, game.CurrentCardCzarId, whiteCardPlay.CardCzarId)
}

// TODO: test czar can select a winner
// TODO: test czar can select a winner and cause a game to end
// TODO: test czar skips black card
