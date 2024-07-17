package gameRepo_test

import (
	"fmt"
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPacksAreLaoded(t *testing.T) {
	t.Parallel()

	if len(gameLogic.AllPacks) == 0 {
		err := gameLogic.LoadPacks()
		assert.NoError(t, err)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	if repo.GameMap == nil {
		t.Log("The game map is nil")
		t.FailNow()
	}
}

func TestCreateGameFail(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	id, _, err := repo.CreateGame(gameLogic.DefaultGameSettings(), "")
	if err == nil {
		t.Log("When a game errors it should not be made")
		t.FailNow()
	}

	_, found := repo.GameMap[id]
	if found {
		t.Log("The game should not be in the map")
		t.FailNow()
	}
}

func TestCreateGame(t *testing.T) {
	repo := gameRepo.New()

	gameSettings := gameLogic.DefaultGameSettings()
	name := "Dave"
	id, pid, err := repo.CreateGame(gameSettings, name)
	if err != nil {
		t.Log("The game should have been made", err)
		t.FailNow()
	}

	game, _ := repo.GetGame(id)
	if game.PlayersMap[game.GameOwnerId].Name != name {
		t.Log("The player was not made with the correct name")
		t.FailNow()
	}

	assert.Equal(t, pid, game.GameOwnerId, "Game owner should be the returned player ID")
}

func TestGetGames(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()

	assert.Equal(t, repo.GetGames(), []*gameLogic.Game{}, "There should be no games in the repo yet")

	gameSettings := gameLogic.DefaultGameSettings()
	name := "Dave"
	id, _, err := repo.CreateGame(gameSettings, name)
	if err != nil {
		t.Log("The game should have been made", err)
		t.FailNow()
	}

	games := repo.GetGames()
	assert.Contains(t, games, repo.GameMap[id], "The game should be in the games returned by the repo")
	assert.Len(t, games, 1, "There should only be one game in the repo")
}

func TestGameChangeSettingsCannotFindGame(t *testing.T) {
	repo := gameRepo.New()
	err := repo.ChangeSettings(uuid.New(), *gameLogic.DefaultGameSettings())
	assert.Error(t, err, "There should be an error when the game does not exist")
}

func TestGameChangeSettings(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	gameSettings := *gameLogic.DefaultGameSettings()
	name := "Dave"
	gid, _, err := repo.CreateGame(&gameSettings, name)
	assert.NoError(t, err)

	newSettings := *&gameSettings
	newSettings.Password = "password123"
	newSettings.MaxPlayers = 7

	assert.True(t, newSettings.Validate())

	err = repo.ChangeSettings(gid, newSettings)
	assert.NoError(t, err, "The settings should have been changed")

	game, err := repo.GetGame(gid)
	assert.NoError(t, err, "The game should exist")

	assert.Equal(t, newSettings, *game.Settings, "The settings should have been changed")
}

func TestGameChangeSettingsInvalidSettings(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	gameSettings := gameLogic.DefaultGameSettings()
	name := "Dave"
	gid, _, err := repo.CreateGame(gameSettings, name)
	assert.NoError(t, err)

	newSettings := *gameSettings
	newSettings.MaxPlayers = 0

	assert.False(t, newSettings.Validate())

	err = repo.ChangeSettings(gid, newSettings)
	assert.Error(t, err, "The settings should not have been changed")

	game, err := repo.GetGame(gid)
	assert.NoError(t, err)
	assert.Equal(t, *gameSettings, *game.Settings, "The game settings should not have changed")
}

func TestPlayerConnect(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	gameSettings := gameLogic.DefaultGameSettings()
	name := "Dave"
	gid, pid, err := repo.CreateGame(gameSettings, name)
	assert.NoError(t, err)

	err = repo.ConnectPlayer(gid, pid)
	assert.NoError(t, err, "The player should have been connected")

	game, err := repo.GetGame(gid)
	assert.NoError(t, err)
	assert.True(t, game.PlayersMap[pid].Connected)
}

func TestPlayerDisconnect(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	gameSettings := gameLogic.DefaultGameSettings()
	name := "Dave"
	gid, pid, err := repo.CreateGame(gameSettings, name)
	assert.NoError(t, err)

	err = repo.ConnectPlayer(gid, pid)
	assert.NoError(t, err, "The player should have been connected")

	game, err := repo.GetGame(gid)
	assert.NoError(t, err)
	assert.True(t, game.PlayersMap[pid].Connected)

	err = repo.DisconnectPlayer(gid, pid)
	assert.NoError(t, err, "The player should have been disconnected")

	game, err = repo.GetGame(gid)
	assert.NoError(t, err)
	assert.False(t, game.PlayersMap[pid].Connected)
}

func TestPlayerLeaveGameNoGame(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	_, err := repo.PlayerLeaveGame(uuid.New(), uuid.New())

	assert.Error(t, err)
}

func TestPlayerLeave(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	settings := gameLogic.DefaultGameSettings()
	gid, _, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	newPid, err := repo.CreatePlayer(gid, "Bill", settings.Password)
	assert.NoError(t, err)

	res, err := repo.PlayerLeaveGame(gid, newPid)
	assert.NoError(t, err)

	assert.Equal(t, 1, res.PlayersLeft)
}

func TestPlayerLeaveLastPlayer(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	settings := gameLogic.DefaultGameSettings()
	gid, pid, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	res, err := repo.PlayerLeaveGame(gid, pid)
	assert.NoError(t, err)
	assert.Equal(t, 0, res.PlayersLeft)

	_, err = repo.GetGame(gid)
	assert.Error(t, err)
}

func TestPlayerPlaysCards(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	settings := gameLogic.DefaultGameSettings()
	gid, pid, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = repo.CreatePlayer(gid, "Player 2", settings.Password)
	assert.NoError(t, err)

	_, err = repo.CreatePlayer(gid, "Player 3", settings.Password)
	assert.NoError(t, err)

	startGameInfo, err := repo.StartGame(gid)
	assert.NoError(t, err)

	cards := make([]int, 0)
	for i := 0; i < int(startGameInfo.CurrentBlackCard.CardsToPlay); i++ {
		cards = append(cards, startGameInfo.PlayerHands[pid][i].Id)
	}

	info, err := repo.PlayerPlayCards(gid, pid, cards)
	assert.NoError(t, err)
	assert.False(t, info.MovedToNextCardCzarPhase)
}

func TestPlayerPlaysCardsNotInGame(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	_, err := repo.PlayerPlayCards(uuid.New(), uuid.New(), []int{1, 2, 3})
	assert.Error(t, err)
}

func TestGamesWithNoPlayersOlderThanMaxEmptyAgeAreRemoved(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	settings := gameLogic.DefaultGameSettings()

	gid, pid, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	game, err := repo.GetGame(gid)
	assert.NoError(t, err)

	game.LastAction = game.LastAction.Add(-(gameRepo.MaxGameWithNoPlayersAge + 1))
	assert.True(t, game.TimeSinceLastAction() > gameRepo.MaxGameWithNoPlayersAge)

	player, found := game.PlayersMap[pid]
	assert.True(t, found)
	player.Connected = false

	gid2, pid, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	game, err = repo.GetGame(gid2)
	assert.NoError(t, err)

	game.LastAction = game.LastAction.Add(-(gameRepo.MaxGameWithNoPlayersAge + 1))
	assert.True(t, game.TimeSinceLastAction() > gameRepo.MaxGameWithNoPlayersAge)

	player, found = game.PlayersMap[pid]
	assert.True(t, found)
	player.Connected = false

	gamesEnded := repo.EndOldGames()
	assert.Len(t, gamesEnded, 2)
	assert.Contains(t, gamesEnded, gid)
	assert.Contains(t, gamesEnded, gid2)

	_, err = repo.GetGame(gid)
	assert.Error(t, err)

	_, err = repo.GetGame(gid2)
	assert.Error(t, err)
}

func TestGamesWithNoPlayersYoungerThanMaxEmptyAgeAreNotRemoved(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	settings := gameLogic.DefaultGameSettings()

	gid, pid, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	game, err := repo.GetGame(gid)
	assert.NoError(t, err)

	player, found := game.PlayersMap[pid]
	assert.True(t, found)
	player.Connected = false

	gid2, pid, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	game, err = repo.GetGame(gid2)
	assert.NoError(t, err)

	player, found = game.PlayersMap[pid]
	assert.True(t, found)
	player.Connected = false

	gamesEnded := repo.EndOldGames()
	assert.Len(t, gamesEnded, 0)

	_, err = repo.GetGame(gid)
	assert.NoError(t, err)

	_, err = repo.GetGame(gid2)
	assert.NoError(t, err)
}

func TestGamesOlderThanMaxAgeAreEnded(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	settings := gameLogic.DefaultGameSettings()

	gid, _, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	game, err := repo.GetGame(gid)
	assert.NoError(t, err)

	game.LastAction = game.LastAction.Add(-(gameRepo.MaxGameInLobbyAge + 1))
	assert.True(t, game.TimeSinceLastAction() > gameRepo.MaxGameInLobbyAge)

	gamesEnded := repo.EndOldGames()
	assert.Len(t, gamesEnded, 1)
	assert.Contains(t, gamesEnded, gid)

	_, err = repo.GetGame(gid)
	assert.Error(t, err)
}

func TestGamesYoungerThanMaxAgeAreNotEnded(t *testing.T) {
	t.Parallel()

	repo := gameRepo.New()
	settings := gameLogic.DefaultGameSettings()

	gid, _, err := repo.CreateGame(settings, "Dave")
	assert.NoError(t, err)

	_, err = repo.GetGame(gid)
	assert.NoError(t, err)

	gamesEnded := repo.EndOldGames()
	assert.Len(t, gamesEnded, 0)

	_, err = repo.GetGame(gid)
	assert.NoError(t, err)
}

var table = []struct {
	input int
}{
	{input: 100},
	{input: 500},
	{input: 1000},
	{input: 5000},
}

func BenchmarkEndGames(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", v.input), func(b *testing.B) {

			b.ReportAllocs()
			b.ReportMetric(float64(b.N), "games")

			repo := gameRepo.New()
			settings := gameLogic.DefaultGameSettings()

			for i := 0; i < b.N; i++ {
				gid, _, err := repo.CreateGame(settings, "Dave")
				if err != nil {
					b.FailNow()
				}

				game, err := repo.GetGame(gid)
				if err != nil {
					b.FailNow()
				}

				game.LastAction = game.LastAction.Add(-(gameRepo.MaxGameInLobbyAge + 1))
			}

			b.ResetTimer()
			repo.EndOldGames()
		})
	}
}
