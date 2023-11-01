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
