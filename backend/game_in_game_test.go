package main

import (
	"fmt"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func (s *ServerTestSuite) TestStartGameTooFewPlayers() {
	t := s.T()
	t.Parallel()

	game := createTestGame(t)
	url := WsBaseUrl + "/games/join"

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	conn, _, err := dialer.Dial(url, game.Jar.Headers())
	require.Nil(t, err, "Should have connected to the ws server successfully")
	defer conn.Close()
	require.NotNil(t, conn)

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := conn.ReadMessage()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Nil(t, err)
	require.Equal(t, game.Ids.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state

	msgType, msg, err = conn.ReadMessage()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	require.Nil(t, err, "Should be a join message")
	require.Equal(t, game.Ids.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        game.Ids.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	require.Nil(t, err)

	conn.WriteMessage(websocket.TextMessage, msgBytes)
	_, msg, err = conn.ReadMessage()
	require.Nil(t, err, "Should be able to read the message")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcCommandErrorMsg](msg)
	require.Nil(t, err, "Should be a command error message")
	require.NotEmpty(t, rpcMsg.Reason)
}

func (s *ServerTestSuite) TestStartGameEnoughPlayers() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Nil(t, err)
	require.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	require.NoError(t, err, "Should be a join message")
	require.Equal(t, client.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	for i := 1; i < gameLogic.MinPlayers; i++ {
		name := fmt.Sprintf("Player %d", i)
		info, err := client.AddPlayer(name)
		require.NoError(t, err)

    s.ReadCreateJoinMessages(t, client, info.PlayerId)
	}

	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	require.Nil(t, err)

	err = client.Write(msgBytes)
	require.NoError(t, err)

	_, msg, err = client.Read()
	require.NoError(t, err, "Should be able to read the message")
	require.NotEmpty(t, msg, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcRoundInformationMsg](msg)
	require.NoError(t, err)
	require.Len(t, rpcMsg.YourHand, gameLogic.HandSize)
	require.NotEmpty(t, rpcMsg.BlackCard)
	require.Empty(t, rpcMsg.YourPlays)
	require.Equal(t, 0, rpcMsg.TotalPlays)
	require.Equal(t, uint(1), rpcMsg.RoundNumber)
	require.Empty(t, rpcMsg.PreviousWinnerDetails)
}

func (s *ServerTestSuite) TestPlayingCardInGame() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Nil(t, err)
	require.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	require.NoError(t, err, "Should be a join message")
	require.Equal(t, client.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	for i := 1; i < gameLogic.MinPlayers; i++ {
		name := fmt.Sprintf("Player %d", i)
		info, err := client.AddPlayer(name)
		require.NoError(t, err)
    require.NotEmpty(t, info)
    s.ReadCreateJoinMessages(t, client, info.PlayerId)
	}

	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	require.NoError(t, err)

	err = client.Write(msgBytes)
	require.NoError(t, err)

	_, msg, err = client.Read()
	require.NoError(t, err, "Should be able to read the message")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcRoundInformationMsg](msg)
	require.NoError(t, err)
	require.Len(t, rpcMsg.YourHand, gameLogic.HandSize)
	require.NotEmpty(t, rpcMsg.BlackCard)
	require.Empty(t, rpcMsg.YourPlays)
	require.Equal(t, 0, rpcMsg.TotalPlays)
	require.Equal(t, uint(1), rpcMsg.RoundNumber)
	require.Empty(t, rpcMsg.PreviousWinnerDetails)

	playCardMsg, err := network.EncodeRpcMessage(network.RpcPlayCardsMsg{CardIds: []int{rpcMsg.YourHand[0].Id}})
	require.NoError(t, err)

	err = client.Write(playCardMsg)
	require.NoError(t, err)

	_, msg, err = client.Read()
	require.NoError(t, err)

	cardPlayedMsg, err := network.DecodeAs[network.RpcOnCardPlayedMsg](msg)
	require.NoError(t, err)
	require.Equal(t, client.PlayerId, cardPlayedMsg.PlayerId)
}

func (s *ServerTestSuite) TestPlayersGetRoundInfoAfterWinnerSelected() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Nil(t, err)
	require.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
  onJoinMsg := ReadMessage[network.RpcOnJoinMsg](s, t, client)
	require.Equal(t, client.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	// Add players
	for i := 1; i < gameLogic.MinPlayers; i++ {
		name := fmt.Sprintf("Player %d", i)
		info, err := client.AddPlayer(name)
		require.NoError(t, err)

    s.ReadCreateJoinMessages(t, client, info.PlayerId)
	}

	// Start game
	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	require.Nil(t, err)

	err = client.Write(msgBytes)
	require.NoError(t, err)

  rpcMsg := ReadMessage[network.RpcRoundInformationMsg](s, t, client)
	require.NoError(t, err)
	require.Len(t, rpcMsg.YourHand, gameLogic.HandSize)
	require.NotEmpty(t, rpcMsg.BlackCard)
	require.Empty(t, rpcMsg.YourPlays)
	require.Equal(t, 0, rpcMsg.TotalPlays)
	require.Equal(t, uint(1), rpcMsg.RoundNumber)
	require.Empty(t, rpcMsg.PreviousWinnerDetails)

	// Play card
	game, err := gameRepo.Repo.GetGame(client.GameId)
	require.NoError(t, err)

	oldBlackCard := game.CurrentBlackCard
	whiteCards := make([]*gameLogic.WhiteCard, 0)
	whiteCardIds := make([]int, 0)
	cardsToPlay := game.CurrentBlackCard.CardsToPlay

	for i := range cardsToPlay {
		whiteCards = append(whiteCards, game.PlayersMap[client.PlayerId].Hand[int(i)])
		whiteCardIds = append(whiteCardIds, int(i))
	}

	playCardMsg, err := network.EncodeRpcMessage(network.RpcPlayCardsMsg{CardIds: whiteCardIds})
	require.NoError(t, err)

	err = client.Write(playCardMsg)
	require.NoError(t, err)

  cardPlayedMsg := ReadMessage[network.RpcOnCardPlayedMsg](s, t, client)
	require.Equal(t, client.PlayerId, cardPlayedMsg.PlayerId)

	// Select the winner
	info, err := gameRepo.Repo.CzarSelectsCard(client.GameId, game.CurrentCardCzarId, whiteCardIds)
	require.NoError(t, err)

	require.Equal(t, info.WinnerId, client.PlayerId)
	require.NotEqual(t, oldBlackCard, info.NewBlackCard)
	require.False(t, info.GameEnded)
	require.Contains(t, game.Players, info.NewCzarId)

	for pid, hand := range info.Hands {
		require.Contains(t, game.Players, pid)
		require.Len(t, hand, 7)

		for _, card := range hand {
			require.NotEmpty(t, card)
		}
	}

	expectedPreviousWinner := gameLogic.PreviousWinner{
		PlayerId:   client.PlayerId,
		BlackCard:  oldBlackCard,
		Whitecards: whiteCards,
	}
	require.Equal(t, expectedPreviousWinner, info.PreviousWinner)
	require.Equal(t, info.PreviousWinner, game.PreviousWinner)

	// The player should have had a RpcOnWhiteCardPlayPhase
	require.NoError(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	whiteCardPlay, err := network.DecodeAs[network.RpcOnWhiteCardPlayPhase](msg)
	require.NoError(t, err, "Should be a white card play message")
	require.Equal(t, expectedPreviousWinner, whiteCardPlay.Winner)
	require.Equal(t, client.PlayerId, whiteCardPlay.WinnerId)
	require.Equal(t, game.CurrentBlackCard, whiteCardPlay.BlackCard)
	require.Len(t, whiteCardPlay.YourHand, 7)
	require.Equal(t, game.CurrentCardCzarId, whiteCardPlay.CardCzarId)
}

func (s *ServerTestSuite) TestMulliganHand() {
	t := s.T()
	t.Parallel()

	client, err := NewTestGameConnection()
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// First message should be the player join broadcast, which we ignore
	msgType, msg, err := client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onPlayerJoinMsg, err := network.DecodeAs[network.RpcOnPlayerJoinMsg](msg)
	require.Nil(t, err)
	require.Equal(t, client.PlayerId, onPlayerJoinMsg.Id, "The current user should have joined the game")

	// Second message should be the state
	msgType, msg, err = client.Read()

	require.Nil(t, err, "Should be able to read (the initial game state)")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")
	require.Equal(t, msgType, websocket.TextMessage)

	onJoinMsg, err := network.DecodeAs[network.RpcOnJoinMsg](msg)
	require.NoError(t, err, "Should be a join message")
	require.Equal(t, client.GameId, onJoinMsg.State.Id)
	require.Len(t, onJoinMsg.State.Players, 1)
	require.Contains(t, onJoinMsg.State.Players, gameLogic.Player{
		Id:        client.PlayerId,
		Name:      "Dave",
		Points:    0,
		Connected: true})

	for i := 1; i < gameLogic.MinPlayers; i++ {
		name := fmt.Sprintf("Player %d", i)
		info, err := client.AddPlayer(name)
		require.NoError(t, err)

    s.ReadCreateJoinMessages(t, client, info.PlayerId)
	}

	startGameMsg := network.RpcStartGameMsg{}
	msgBytes, err := network.EncodeRpcMessage(startGameMsg)
	require.Nil(t, err)

	err = client.Write(msgBytes)
	require.NoError(t, err)

	_, msg, err = client.Read()
	require.Nil(t, err, "Should be able to read the message")
	require.True(t, len(msg) > 0, "Message should have a non-zero length")

	rpcMsg, err := network.DecodeAs[network.RpcRoundInformationMsg](msg)
	require.NoError(t, err)
	require.Len(t, rpcMsg.YourHand, gameLogic.HandSize)
	require.NotEmpty(t, rpcMsg.BlackCard)
	require.Empty(t, rpcMsg.YourPlays)
	require.Equal(t, 0, rpcMsg.TotalPlays)
	require.Equal(t, uint(1), rpcMsg.RoundNumber)
	require.Empty(t, rpcMsg.PreviousWinnerDetails)

	mulliganMsg, err := network.EncodeRpcMessage(network.RpcMulliganHand{})
	require.NoError(t, err)

	err = client.Write(mulliganMsg)
	require.NoError(t, err)

	_, msg, err = client.Read()
	require.NoError(t, err)

	cardPlayedMsg, err := network.DecodeAs[network.RpcOnNewHand](msg)
	require.NoError(t, err)
	require.Len(t, cardPlayedMsg.WhiteCards, 7)
}

// TODO: test czar can select a winner
// TODO: test czar can select a winner and cause a game to end
// TODO: test czar skips black card
