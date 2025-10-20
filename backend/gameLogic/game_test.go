package gameLogic_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultGameSettings(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()

	if settings == nil {
		t.Log("Default game settings are nil")
		t.FailNow()
	}

	if !settings.Validate() {
		t.Log("Default game settings should be valid")
		t.FailNow()
	}
}

func TestGameSettingsValidateMinRounds(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	settings.MaxRounds = gameLogic.MinRounds - 1
	if settings.Validate() {
		t.Log("Values below minimum rounds should not validate")
		t.FailNow()
	}
}

func TestGameSettingsValidateMaxRounds(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	settings.MaxRounds = gameLogic.MaxRounds + 1
	if settings.Validate() {
		t.Log("Values above maximum ronuds should not validate")
		t.FailNow()
	}
}

func TestGameSettingsValidateMinPlayingToPoints(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	settings.PlayingToPoints = gameLogic.MinPlayingToPoints - 1
	if settings.Validate() {
		t.Log("Values below minimum playing to points should not validate")
		t.FailNow()
	}
}

func TestGameSettingsValidateMaxPlayingToPoints(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	settings.PlayingToPoints = gameLogic.MaxPlayingToPoints + 1
	if settings.Validate() {
		t.Log("Values above maximum playing to points should not validate")
		t.FailNow()
	}
}

func TestGameSettingsMaximumPasswordLength(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	settings.Password = strings.Repeat("a", gameLogic.MaxPasswordLength+1)
	if settings.Validate() {
		t.Log("Passwords above maximum length should fail")
		t.FailNow()
	}
}

func TestNewGameInvalidSettings(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	settings.MaxPlayers = gameLogic.MinPlayers - 1
	if settings.Validate() {
		t.Log("The invalid settings are valid")
		t.FailNow()
	}

	name := "Dave"
	_, err := gameLogic.NewGame(settings, name)
	if err == nil {
		t.Log(fmt.Sprintf("The error should be nil %s", err))
		t.FailNow()
	}
}

func TestNewGameInvalidPlayer(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	if !settings.Validate() {
		t.Log("The valid settings are invalid")
		t.FailNow()
	}

	name := ""
	_, err := gameLogic.NewGame(settings, name)
	if err == nil {
		t.Log(fmt.Sprintf("The error should be nil %s", err))
		t.FailNow()
	}
}

func TestNewGame(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	if !settings.Validate() {
		t.Log("The settings should be valid")
		t.FailNow()
	}

	name := "Dave"
	game, err := gameLogic.NewGame(settings, name)
	assert.NoError(t, err)
	assert.NotNil(t, game)

	assert.Equal(t, settings, game.Settings)
	assert.NotEmpty(t, game.CreationTime)
	assert.Empty(t, game.CurrentCardCzarId)
	assert.NotNil(t, game.Players)
	assert.NotEmpty(t, game.GameOwnerId)

	_, found := game.PlayersMap[game.GameOwnerId]
	assert.True(t, found)
	assert.Len(t, game.Players, 1)
	assert.Equal(t, game.Players[0], game.GameOwnerId)
	assert.Equal(t, uint(0), game.CurrentRound)
	assert.NotEmpty(t, game.Id)
	assert.NotEmpty(t, game.LastAction)
	assert.NotEmpty(t, game.TimeSinceLastAction())
  assert.Empty(t, game.PreviousWinner)
}

func TestGameInfo(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.Nil(t, err, "There should not be an error with making the game", err)

	oldTime := game.LastAction
	info := game.Info()
	assert.Equal(t, oldTime, game.LastAction)

	assert.Equal(t, game.Id, info.Id, "Game IDs should be equal")
	assert.Equal(t, 1, info.PlayerCount, "There should only be one player")
	assert.Equal(t, game.Settings.MaxPlayers, info.MaxPlayers, "Max players should be equal")
	assert.False(t, info.HasPassword, "The game should not be marked as having a password")

	game.Settings.Password = "poop"
	info = game.Info()

	assert.True(t, info.HasPassword, "The game should be marked as having a password")
}

func TestGameStateInfo(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.Nil(t, err, "There should not be an error with making the game", err)

	pid := game.Players[0]
	info := game.StateInfo(pid)

	expectedPlayers := make([]gameLogic.Player, len(game.Players))
	for i, pid := range game.Players {
		expectedPlayers[i] = gameLogic.Player{
			Id:        pid,
			Name:      game.PlayersMap[pid].Name,
			Connected: false,
			Points:    0,
		}
	}

	assert.Equal(t, game.Id, info.Id)
	assert.Equal(t, *game.Settings, info.Settings)
	assert.Equal(t, expectedPlayers, info.Players)
	assert.Equal(t, game.CreationTime, info.CreationTime)
	assert.Equal(t, game.GameState, info.GameState)
	assert.Equal(t, game.GameOwnerId, info.GameOwnerId)
	assert.Len(t, info.RoundInfo.YourHand, 0)
	assert.Len(t, info.RoundInfo.PlayersWhoHavePlayed, 0)
	assert.Empty(t, info.RoundInfo.BlackCard)
	assert.Empty(t, info.RoundInfo.CardCzarId)
	assert.Empty(t, info.RoundInfo.RoundNumber)
  assert.Empty(t, info.RoundInfo.PreviousWinner)
}

func TestGameStateInfoMidRound(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	game.GameState = gameLogic.GameStateWhiteCardsBeingSelected
	assert.Nil(t, err, "There should not be an error with making the game", err)

	pid := game.Players[0]
	info := game.StateInfo(pid)

	expectedPlayers := make([]gameLogic.Player, len(game.Players))
	for i, pid := range game.Players {
		expectedPlayers[i] = gameLogic.Player{
			Id:        pid,
			Name:      game.PlayersMap[pid].Name,
			Connected: false,
			Points:    0,
		}
	}

	assert.Equal(t, game.Id, info.Id)
	assert.Equal(t, *game.Settings, info.Settings)
	assert.Equal(t, expectedPlayers, info.Players)
	assert.Equal(t, game.CreationTime, info.CreationTime)
	assert.Equal(t, game.GameState, info.GameState)
	assert.Equal(t, game.GameOwnerId, info.GameOwnerId)
  assert.Equal(t, game.PreviousWinner, info.RoundInfo.PreviousWinner)
}

func TestAddInvalidPlayer(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	_, err = game.AddPlayer("")
	if err == nil {
		t.Log("Should not be able to add an invalid player")
		t.FailNow()
	}
}

func TestAddingDuplicatePlayer(t *testing.T) {
	t.Parallel()

	name := "Dave"
	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, name)
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	_, err = game.AddPlayer(name)
	if err == nil {
		t.Log("Should not be able to add an invalid player")
		t.FailNow()
	}
}

func TestAddPlayer(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	id, err := game.AddPlayer("Steve")
	if err != nil {
		t.Log("Should not be able to add an invalid player")
		t.FailNow()
	}

	var nilUuid uuid.UUID
	if id == nilUuid {
		t.Log("Id was not set")
		t.FailNow()
	}

	player, found := game.PlayersMap[id]
	if !found {
		t.Log("Cannot find the player in the player map")
		t.FailNow()
	}

	if player.Id != id {
		t.Log("Player ID mismatch between list and map")
		t.FailNow()
	}

	if len(game.Players) != 2 {
		t.Log("There should be two players")
		t.FailNow()
	}

	if game.Players[len(game.Players)-1] != id {
		t.Log("Player not at end of players list")
		t.FailNow()
	}
}

func TestAddingPlayers(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	for i := 0; i < int(settings.MaxPlayers)-1; i++ {
		id, err := game.AddPlayer(fmt.Sprintf("Steve %d", i))
		if err != nil {
			t.Log("Should not be able to add an invalid player")
			t.FailNow()
		}

		var nilUuid uuid.UUID
		if id == nilUuid {
			t.Log("Id was not set")
			t.FailNow()
		}

		player, found := game.PlayersMap[id]
		if !found {
			t.Log("Cannot find the player in the player map")
			t.FailNow()
		}

		if player.Id != id {
			t.Log("Player ID mismatch between list and map")
			t.FailNow()
		}

		if len(game.Players) != 2+i {
			t.Log("There should be two players")
			t.FailNow()
		}

		if game.Players[len(game.Players)-1] != id {
			t.Log("Player not at end of players list")
			t.FailNow()
		}
	}
}

func TestAddingMaxPlayers(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	for i := 0; i < int(settings.MaxPlayers)-1; i++ {
		id, err := game.AddPlayer(fmt.Sprintf("Steve %d", i))
		if err != nil {
			t.Log("Should not be able to add an invalid player")
			t.FailNow()
		}

		var nilUuid uuid.UUID
		if id == nilUuid {
			t.Log("Id was not set")
			t.FailNow()
		}

		player, found := game.PlayersMap[id]
		if !found {
			t.Log("Cannot find the player in the player map")
			t.FailNow()
		}

		if player.Id != id {
			t.Log("Player ID mismatch between list and map")
			t.FailNow()
		}

		if len(game.Players) != 2+i {
			t.Log("There should be two players")
			t.FailNow()
		}

		if game.Players[len(game.Players)-1] != id {
			t.Log("Player not at end of players list")
			t.FailNow()
		}
	}

	_, err = game.AddPlayer(fmt.Sprintf("Final Steve"))
	if err == nil {
		t.Log("Should not be able to add beyond the maximum amount of players")
		t.FailNow()
	}
}

func TestStartGameOnePlayerFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	_, err = game.StartGame()
	if err == nil {
		t.Log("Should not be able to start a game with 1 players")
		t.FailNow()
	}
}

func TestStartGameLessThanMinPlayerFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	for i := 0; i < gameLogic.MinPlayers-2; i++ {
		_, err = game.AddPlayer(fmt.Sprintf("Player %d", i))
		if err != nil {
			t.Log("Cannot add a player", err)
			t.FailNow()
		}
	}

	_, err = game.StartGame()
	if err == nil {
		t.Log("Should not be able to start a game with min players")
		t.FailNow()
	}
}

func TestStartGameNotInLobbyFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	game.GameState = gameLogic.GameStateCzarJudgingCards

	_, err = game.StartGame()
	if err == nil {
		t.Log("Should not be able to start a game with invalid state")
		t.FailNow()
	}
}

func TestCreateGameWithNoPacksFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	settings.CardPacks = []uuid.UUID{}
	_, err := gameLogic.NewGame(settings, "Dave")
	assert.Error(t, err)
}

func TestStartGameSuccess(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	for i := 0; i < int(settings.MaxPlayers)-1; i++ {
		_, err = game.AddPlayer(fmt.Sprintf("Player %d", i))
		if err != nil {
			t.Log("Cannot add a player", err)
			t.FailNow()
		}
	}

	info, err := game.StartGame()
	assert.NoError(t, err)

	for _, player := range game.PlayersMap {
		assert.Equal(t, gameLogic.HandSize, player.CardsInHand())
		assert.Len(t, info.PlayerHands[player.Id], gameLogic.HandSize)
		for _, card := range player.Hand {
			assert.Contains(t, info.PlayerHands[player.Id], card)
		}

		assert.Len(t, info.PlayersPlays[player.Id], 0)
	}

	assert.Equal(t, info.CurrentBlackCard, game.CurrentBlackCard)
	assert.Equal(t, info.CurrentCardCzarId, game.CurrentCardCzarId)
	assert.Equal(t, info.RoundNumber, game.CurrentRound)
	assert.NotEmpty(t, info.CurrentBlackCard)
	assert.NotEmpty(t, info.CurrentCardCzarId)
}

func TestAddingPlayerToGameInProgress(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	for i := 0; i < int(settings.MaxPlayers)-2; i++ {
		_, err = game.AddPlayer(fmt.Sprintf("Player %d", i))
		if err != nil {
			t.Log("Cannot add a player", err)
			t.FailNow()
		}
	}

	_, err = game.StartGame()
	assert.NoError(t, err)

	pid, err := game.AddPlayer("Jesus Christ")
	assert.NoError(t, err)

	player := game.PlayersMap[pid]
	assert.Len(t, player.Hand, gameLogic.HandSize)
	assert.Len(t, player.CurrentPlay, 0)
}

func TestRemovePlayerFromGameThatIsNotInThere(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.RemovePlayer(uuid.New())
	assert.Error(t, err, "Should not be able to remove a player that does not exist")

	assert.Len(t, game.Players, 1)
	assert.Len(t, game.PlayersMap, 1)
}

func TestRemovingLastPlayer(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	res, err := game.RemovePlayer(game.Players[0])
	assert.NoError(t, err)
	assert.Equal(t, res.PlayersLeft, 0)

	assert.Len(t, game.Players, 0)
	assert.Len(t, game.PlayersMap, 0)
}

func TestRemovingGameOwnerReassignsIt(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	pid, err := game.AddPlayer("John")
	assert.NoError(t, err)

	res, err := game.RemovePlayer(game.GameOwnerId)
	assert.NoError(t, err)

	assert.Equal(t, game.GameOwnerId, res.NewGameOwner)
	assert.Equal(t, pid, game.GameOwnerId)
	assert.Equal(t, res.PlayersLeft, 1)

	assert.Len(t, game.Players, 1)
	assert.Len(t, game.PlayersMap, 1)
}

func TestChangeSettingsValid(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	newSettings := gameLogic.DefaultGameSettings()
	newSettings.MaxPlayers = 10

	err = game.ChangeSettings(*newSettings)
	assert.NoError(t, err)
}

func TestCannotChangeSettingsInLobby(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	game.GameState = gameLogic.GameStateWhiteCardsBeingSelected
	newSettings := gameLogic.DefaultGameSettings()
	newSettings.MaxPlayers = 10

	err = game.ChangeSettings(*newSettings)
	assert.Error(t, err)
}

func TestCannotChangeSettingsToInvalidSettings(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	newSettings := gameLogic.DefaultGameSettings()
	newSettings.MaxPlayers = 0

	err = game.ChangeSettings(*newSettings)
	assert.Error(t, err)
}

func TestPlayingCardThatDoesNotExistFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	game.GameState = gameLogic.GameStateWhiteCardsBeingSelected
	game.CurrentBlackCard = &gameLogic.BlackCard{
		Id:          180,
		CardsToPlay: 1,
		BodyText:    "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
	}

	pid := game.Players[0]
	cardId := -1
	_, err = gameLogic.GetWhiteCard(cardId)
	assert.Error(t, err)

	_, err = game.PlayCards(pid, []int{cardId})
	assert.Error(t, err)
}

func TestPlayingCardsForAPlayerThatDoesNotExistFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	game.GameState = gameLogic.GameStateWhiteCardsBeingSelected
	game.CurrentBlackCard = &gameLogic.BlackCard{
		Id:          180,
		CardsToPlay: 1,
		BodyText:    "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
	}

	cardId := 1
	_, err = gameLogic.GetWhiteCard(cardId)
	assert.NoError(t, err)

	_, err = game.PlayCards(uuid.New(), []int{cardId})
	assert.Error(t, err)
}

func TestPlayingWrongAmountOfCardsFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	game.GameState = gameLogic.GameStateWhiteCardsBeingSelected
	game.CurrentBlackCard = &gameLogic.BlackCard{
		Id:          180,
		CardsToPlay: 2,
		BodyText:    "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
	}

	cardId := 1
	card, err := gameLogic.GetWhiteCard(cardId)
	assert.NoError(t, err)

	pid := game.Players[0]
	game.PlayersMap[pid].Hand = make(map[int]*gameLogic.WhiteCard)
	game.PlayersMap[pid].Hand[cardId] = card

	_, err = game.PlayCards(pid, []int{cardId})
	assert.Error(t, err)
	assert.Nil(t, game.PlayersMap[pid].CurrentPlay)
}

func TestPlayingDuplicateCardsFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	game.GameState = gameLogic.GameStateWhiteCardsBeingSelected
	game.CurrentBlackCard = &gameLogic.BlackCard{
		Id:          180,
		CardsToPlay: 2,
		BodyText:    "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
	}

	cardId := 1
	card, err := gameLogic.GetWhiteCard(cardId)
	assert.NoError(t, err)

	pid := game.Players[0]
	game.PlayersMap[pid].Hand = make(map[int]*gameLogic.WhiteCard)
	game.PlayersMap[pid].Hand[cardId] = card

	_, err = game.PlayCards(pid, []int{cardId, cardId})
	assert.Error(t, err)
	assert.Nil(t, game.PlayersMap[pid].CurrentPlay)
}

func TestPlayingCardNotInHandFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	game.GameState = gameLogic.GameStateWhiteCardsBeingSelected
	game.CurrentBlackCard = &gameLogic.BlackCard{
		Id:          180,
		CardsToPlay: 1,
		BodyText:    "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
	}

	cardId := 2
	_, err = gameLogic.GetWhiteCard(cardId)
	assert.NoError(t, err)

	pid := game.Players[0]
	game.PlayersMap[pid].Hand = make(map[int]*gameLogic.WhiteCard)
	game.PlayersMap[pid].Hand[cardId-1] = gameLogic.NewWhiteCard(cardId-1, "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.")

	_, err = game.PlayCards(pid, []int{cardId})
	assert.Error(t, err)
	assert.Nil(t, game.PlayersMap[pid].CurrentPlay)
}

func TestPlayingCardSuccessCase(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 3")
	assert.NoError(t, err)

	_, err = game.StartGame()
	assert.NoError(t, err)

	cardId := 1
	card, err := gameLogic.GetWhiteCard(cardId)
	assert.NoError(t, err)

	pid := game.Players[0]
	game.PlayersMap[pid].Hand = make(map[int]*gameLogic.WhiteCard)
	game.PlayersMap[pid].Hand[cardId] = card

	resp, err := game.PlayCards(pid, []int{cardId})
	assert.NoError(t, err)
	assert.Equal(t, game.PlayersMap[pid].CurrentPlay, []*gameLogic.WhiteCard{game.PlayersMap[pid].Hand[cardId]})
	assert.False(t, resp.MovedToNextCardCzarPhase)
}

func TestPlayingCardCausesCzarJudingPhase(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 3")
	assert.NoError(t, err)

	_, err = game.StartGame()
	assert.NoError(t, err)

	cardId := 1
	card, err := gameLogic.GetWhiteCard(cardId)
	assert.NoError(t, err)

	for _, player := range game.PlayersMap {
		player.Hand = make(map[int]*gameLogic.WhiteCard)
		player.Hand[cardId] = card
	}

	pid := game.Players[0]
	resp, err := game.PlayCards(pid, []int{cardId})
	assert.NoError(t, err)
	assert.Equal(t, game.PlayersMap[pid].CurrentPlay, []*gameLogic.WhiteCard{game.PlayersMap[pid].Hand[cardId]})
	assert.False(t, resp.MovedToNextCardCzarPhase)

	// Check czar cannot play
	pid = game.Players[2]
	resp, err = game.PlayCards(pid, []int{cardId})
	assert.Error(t, err)
	assert.Nil(t, game.PlayersMap[pid].CurrentPlay)
	assert.False(t, resp.MovedToNextCardCzarPhase)

	// Play final play
	pid = game.Players[1]
	resp, err = game.PlayCards(pid, []int{cardId})
	assert.NoError(t, err)

	// Game should have continued
	assert.True(t, resp.MovedToNextCardCzarPhase)

	for _, plays := range resp.CzarJudingPhaseInfo.AllPlays {
		assert.Equal(t, plays, []*gameLogic.WhiteCard{card})
	}
	assert.NotNil(t, resp.CzarJudingPhaseInfo.PlayerHands)

	for _, hand := range resp.CzarJudingPhaseInfo.PlayerHands.Hands {
		assert.NotNil(t, hand)
		assert.Len(t, hand, 7)
	}
}

func TestPlayingCardInWrongGameStateFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	game.GameState = gameLogic.GameStateCzarJudgingCards
	game.CurrentBlackCard = &gameLogic.BlackCard{
		Id:          180,
		CardsToPlay: 1,
		BodyText:    "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
	}

	cardId := 2
	_, err = gameLogic.GetWhiteCard(cardId)
	assert.NoError(t, err)

	pid := game.Players[0]
	game.PlayersMap[pid].Hand = make(map[int]*gameLogic.WhiteCard)
	game.PlayersMap[pid].Hand[cardId-1] = gameLogic.NewWhiteCard(cardId-1, "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.")

	_, err = game.PlayCards(pid, []int{cardId})
	assert.Error(t, err)
	assert.Nil(t, game.PlayersMap[pid].CurrentPlay)
}

// TODO: Check that when a player leaves and all other players have played the judging starts
// TODO: Check that when a player leaves and there are too few players the game ends

func TestJudgingInWrongPhaseFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.CzarSelectCards(game.Players[0], []int{0, 1})
	assert.Error(t, err)
}

func TestJudgingWithPlayerWhoIsNotCzarFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	assert.NoError(t, err)

	startGameInfo, err := game.StartGame()
	assert.NoError(t, err)

	var lastPlay []int
	for i, pid := range game.Players {
		if pid == game.CurrentCardCzarId {
			continue
		}

		cards := make([]int, 0)
		for i := 0; i < int(game.CurrentBlackCard.CardsToPlay); i++ {
			cards = append(cards, startGameInfo.PlayerHands[pid][i].Id)
		}

		lastPlay = cards
		res, err := game.PlayCards(pid, cards)
		assert.NoError(t, err)

		if i == len(game.Players)-1 {
			assert.True(t, res.MovedToNextCardCzarPhase)
		}
	}

	assert.NotEqual(t, game.CurrentCardCzarId, game.Players[0])
	_, err = game.CzarSelectCards(game.Players[0], lastPlay)
	assert.Error(t, err)
}

func TestJudgingWithPlayerInvalidPlaysFails(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	assert.NoError(t, err)

	startGameInfo, err := game.StartGame()
	assert.NoError(t, err)

	var lastPlay []int
	for i, pid := range game.Players {
		if pid == game.CurrentCardCzarId {
			continue
		}

		cards := make([]int, 0)
		for i := 0; i < int(game.CurrentBlackCard.CardsToPlay); i++ {
			cards = append(cards, startGameInfo.PlayerHands[pid][i].Id)
		}

		lastPlay = cards
		res, err := game.PlayCards(pid, cards)
		assert.NoError(t, err)

		if i == len(game.Players)-1 {
			assert.True(t, res.MovedToNextCardCzarPhase)
		}
	}

	assert.Equal(t, gameLogic.GameStateCzarJudgingCards, game.GameState)

	lastPlay[0] += 1
	_, err = game.CzarSelectCards(game.CurrentCardCzarId, lastPlay)
	assert.Error(t, err)
}

func TestJudgingSuccess(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	assert.NoError(t, err)

	startGameInfo, err := game.StartGame()
	assert.NoError(t, err)

	var lastPlay []int
	var winnerId uuid.UUID
	for i, pid := range game.Players {
		if pid == game.CurrentCardCzarId {
			continue
		}

		cards := make([]int, 0)
		for i := 0; i < int(game.CurrentBlackCard.CardsToPlay); i++ {
			cards = append(cards, startGameInfo.PlayerHands[pid][i].Id)
		}

		lastPlay = cards
		winnerId = pid

		res, err := game.PlayCards(pid, cards)
		assert.NoError(t, err)

		if i == len(game.Players)-1 {
			assert.True(t, res.MovedToNextCardCzarPhase)
		}
	}

	assert.Equal(t, gameLogic.GameStateCzarJudgingCards, game.GameState)

	oldCzar := game.CurrentCardCzarId
	oldBlackCard := game.CurrentBlackCard

	res, err := game.CzarSelectCards(oldCzar, lastPlay)
	assert.NoError(t, err)

	assert.False(t, res.GameEnded)

	assert.NotEqual(t, oldCzar, game.CurrentCardCzarId)
	assert.NotEmpty(t, game.CurrentCardCzarId)
	assert.Equal(t, game.CurrentCardCzarId, res.NewCzarId)

	assert.NotEqual(t, oldBlackCard, game.CurrentBlackCard)
	assert.NotEmpty(t, game.CurrentBlackCard)
	assert.Equal(t, game.CurrentBlackCard, res.NewBlackCard)

	assert.Equal(t, res.WinnerId, winnerId)

	assert.Equal(t, gameLogic.GameStateWhiteCardsBeingSelected, game.GameState)
	assert.Equal(t, 1, game.PlayersMap[winnerId].Points)
	assert.Equal(t, uint(2), game.CurrentRound)

	for _, player := range game.PlayersMap {
		assert.Contains(t, res.Hands, player.Id)
		for _, card := range player.Hand {
			assert.Contains(t, res.Hands[player.Id], card)
		}
	}
}

func TestJudgingGameEndMaxRounds(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	assert.NoError(t, err)

	startGameInfo, err := game.StartGame()
	assert.NoError(t, err)

	var lastPlay []int
	var winnerId uuid.UUID
	for i, pid := range game.Players {
		if pid == game.CurrentCardCzarId {
			continue
		}

		cards := make([]int, 0)
		for i := 0; i < int(game.CurrentBlackCard.CardsToPlay); i++ {
			cards = append(cards, startGameInfo.PlayerHands[pid][i].Id)
		}

		lastPlay = cards
		winnerId = pid

		res, err := game.PlayCards(pid, cards)
		assert.NoError(t, err)

		if i == len(game.Players)-1 {
			assert.True(t, res.MovedToNextCardCzarPhase)
		}
	}

	assert.Equal(t, gameLogic.GameStateCzarJudgingCards, game.GameState)

	oldCzar := game.CurrentCardCzarId

	game.CurrentRound = game.Settings.MaxRounds - 1

	res, err := game.CzarSelectCards(oldCzar, lastPlay)
	assert.NoError(t, err)

	assert.True(t, res.GameEnded)
	assert.Equal(t, winnerId, res.WinnerId)
}

func TestJudgingGameEndMaxScore(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	assert.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	assert.NoError(t, err)

	startGameInfo, err := game.StartGame()
	assert.NoError(t, err)

	var lastPlay []int
	var winnerId uuid.UUID
	for i, pid := range game.Players {
		if pid == game.CurrentCardCzarId {
			continue
		}

		cards := make([]int, 0)
		for i := 0; i < int(game.CurrentBlackCard.CardsToPlay); i++ {
			cards = append(cards, startGameInfo.PlayerHands[pid][i].Id)
		}

		lastPlay = cards
		winnerId = pid

		res, err := game.PlayCards(pid, cards)
		assert.NoError(t, err)

		if i == len(game.Players)-1 {
			assert.True(t, res.MovedToNextCardCzarPhase)
		}
	}

	assert.Equal(t, gameLogic.GameStateCzarJudgingCards, game.GameState)

	oldCzar := game.CurrentCardCzarId

	game.PlayersMap[winnerId].Points = int(game.Settings.PlayingToPoints) - 1

	res, err := game.CzarSelectCards(oldCzar, lastPlay)
	assert.NoError(t, err)

	assert.True(t, res.GameEnded)
	assert.Equal(t, winnerId, res.WinnerId)
}

func TestSkipBlackCardWrongState(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	require.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	require.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	require.NoError(t, err)

	_, err = game.SkipBlackCard(game.CurrentCardCzarId)
	require.Error(t, err)
}

func TestSkipBlackCardNotCzar(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	require.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	require.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	require.NoError(t, err)

	_, err = game.StartGame()
	require.NoError(t, err)

	_, err = game.SkipBlackCard(uuid.New())
	require.Error(t, err)
}

func TestSkipBlackCard(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	require.NoError(t, err)

	_, err = game.AddPlayer("Player 1")
	require.NoError(t, err)

	_, err = game.AddPlayer("Player 2")
	require.NoError(t, err)

	_, err = game.AddPlayer("Player 3")
	require.NoError(t, err)

	_, err = game.StartGame()
	require.NoError(t, err)

	for _, pid := range game.Players {
		if pid == game.CurrentCardCzarId {
			continue
		}

		player, found := game.PlayersMap[pid]
		require.True(t, found)

		hand := make([]*gameLogic.WhiteCard, 0)
		for _, card := range player.Hand {
			hand = append(hand, card)
		}

		cards := make([]int, 0)
		for i := range game.CurrentBlackCard.CardsToPlay {
			cards = append(cards, hand[i].Id)
		}

		_, err = game.PlayCards(pid, cards)
		require.NoError(t, err)

		// Only play one of the players cards
		break
	}

	newCard, err := game.SkipBlackCard(game.CurrentCardCzarId)
	require.NoError(t, err)
	require.NotNil(t, newCard)

	for _, pid := range game.Players {
		player, found := game.PlayersMap[pid]
		require.True(t, found)

		require.Len(t, player.CurrentPlay, 0)
		require.Len(t, player.Hand, gameLogic.HandSize)
	}
}
