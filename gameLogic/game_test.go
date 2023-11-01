package gameLogic_test

import (
	"strings"
	"testing"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
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
