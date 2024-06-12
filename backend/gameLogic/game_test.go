package gameLogic_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	if err != nil {
		t.Log(fmt.Sprintf("The game should be valid and created %s", err))
		t.FailNow()
	}

	if game == nil {
		t.Log("Game should not be nil")
		t.FailNow()
	}

	if game.Settings != settings {
		t.Log("Setting in the game were not set, this is bad")
		t.FailNow()
	}

	var nilTime time.Time
	if game.CreationTime == nilTime {
		t.Log("The time was not set")
		t.FailNow()
	}

	var nilUuid uuid.UUID
	if game.CurrentCardCzarId != nilUuid {
		t.Log("The czar should not be set yet")
		t.FailNow()
	}

	if game.Players == nil {
		t.Log("Player map not set")
		t.FailNow()
	}

	if game.GameOwnerId == nilUuid {
		t.Log("Game owner not set")
		t.FailNow()
	}

	_, found := game.PlayersMap[game.GameOwnerId]
	if !found {
		t.Log("The owner is not in the game")
		t.FailNow()
	}

	if len(game.Players) != 1 {
		t.Log("There should be a player in the list")
		t.FailNow()
	}

	if game.Players[0] != game.GameOwnerId {
		t.Log("Player is not in the player turn list")
		t.FailNow()
	}

	if game.CurrentRound != 0 {
		t.Log("The current round should be set to 0")
		t.FailNow()
	}

	if nilUuid == game.Id {
		t.Log("Game ID was not set")
		t.FailNow()
	}
}

func TestGameInfo(t *testing.T) {
	t.Parallel()

	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	assert.Nil(t, err, "There should not be an error with making the game", err)

	info := game.Info()

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

	info := game.StateInfo()

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

	pid = game.Players[1]
	resp, err = game.PlayCards(pid, []int{cardId})
	assert.NoError(t, err)
	assert.Equal(t, game.PlayersMap[pid].CurrentPlay, []*gameLogic.WhiteCard{game.PlayersMap[pid].Hand[cardId]})
	assert.False(t, resp.MovedToNextCardCzarPhase)

	pid = game.Players[2]
	resp, err = game.PlayCards(pid, []int{cardId})
	assert.NoError(t, err)
	assert.NotNil(t, game.PlayersMap[pid].CurrentPlay)
	assert.Empty(t, game.PlayersMap[pid].CurrentPlay)
	assert.True(t, resp.MovedToNextCardCzarPhase)

	assert.Equal(t, resp.CzarJudingPhaseInfo.AllPlays, []*gameLogic.WhiteCard{card, card, card})
	assert.NotNil(t, resp.CzarJudingPhaseInfo.PlayerHands)

	for _, hand := range resp.CzarJudingPhaseInfo.PlayerHands {
		assert.NotNil(t, hand)
	}
}

// TODO: Check that being in the wrong state fails
// TODO: Check that when a player leaves and all other players have played the judging starts
// TODO: Check that when a player leaves and there are too few players the game ends
