package gameLogic_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/google/uuid"
)

func TestDefaultGameSettings(t *testing.T) {
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
	settings := gameLogic.DefaultGameSettings()
	settings.MaxRounds = gameLogic.MinRounds - 1
	if settings.Validate() {
		t.Log("Values below minimum rounds should not validate")
		t.FailNow()
	}
}

func TestGameSettingsValidateMaxRounds(t *testing.T) {
	settings := gameLogic.DefaultGameSettings()
	settings.MaxRounds = gameLogic.MaxRounds + 1
	if settings.Validate() {
		t.Log("Values above maximum ronuds should not validate")
		t.FailNow()
	}
}

func TestGameSettingsValidateMinPlayingToPoints(t *testing.T) {
	settings := gameLogic.DefaultGameSettings()
	settings.PlayingToPoints = gameLogic.MinPlayingToPoints - 1
	if settings.Validate() {
		t.Log("Values below minimum playing to points should not validate")
		t.FailNow()
	}
}

func TestGameSettingsValidateMaxPlayingToPoints(t *testing.T) {
	settings := gameLogic.DefaultGameSettings()
	settings.PlayingToPoints = gameLogic.MaxPlayingToPoints + 1
	if settings.Validate() {
		t.Log("Values above maximum playing to points should not validate")
		t.FailNow()
	}
}

func TestGameSettingsMaximumPasswordLength(t *testing.T) {
	settings := gameLogic.DefaultGameSettings()
	settings.Password = strings.Repeat("a", gameLogic.MaxPasswordLength+1)
	if settings.Validate() {
		t.Log("Passwords above maximum length should fail")
		t.FailNow()
	}
}

func TestNewGameInvalidSettings(t *testing.T) {
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

func TestAddInvalidPlayer(t *testing.T) {
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
	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	err = game.StartGame()
	if err == nil {
		t.Log("Should not be able to start a game with 1 players")
		t.FailNow()
	}
}

func TestStartGameLessThanMinPlayerFails(t *testing.T) {
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

	err = game.StartGame()
	if err == nil {
		t.Log("Should not be able to start a game with min players")
		t.FailNow()
	}
}

func TestStartGameNotInLobbyFails(t *testing.T) {
	settings := gameLogic.DefaultGameSettings()
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	game.GameState = gameLogic.GameStateCzarJudgingCards

	err = game.StartGame()
	if err == nil {
		t.Log("Should not be able to start a game with invalid state")
		t.FailNow()
	}
}

func TestStartGameNoCardsInDeckFails(t *testing.T) {
	settings := gameLogic.DefaultGameSettings()
	settings.CardPacks = []*gameLogic.CardPack{{}}
	game, err := gameLogic.NewGame(settings, "Dave")
	if err != nil {
		t.Log("Cannot make the game")
		t.FailNow()
	}

	game.GameState = gameLogic.GameStateCzarJudgingCards

	err = game.StartGame()
	if err == nil {
		t.Log("Should not be able to start a game with no cards in the deck")
		t.FailNow()
	}
}
