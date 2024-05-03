package gameRepo

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
	"github.com/google/uuid"
)

const (
	MaxGameInProgressAge    = time.Hour * 3
	MaxGameInLobbyAge       = time.Minute * 15
	MaxGameWithNoPlayersAge = time.Second * 2
)

type GameRepo struct {
	GameMap    map[uuid.UUID]*gameLogic.Game
	GameAgeMap map[uuid.UUID]time.Time
	lock       sync.RWMutex
}

func New() *GameRepo {
	return &GameRepo{GameMap: make(map[uuid.UUID]*gameLogic.Game),
		GameAgeMap: make(map[uuid.UUID]time.Time)}
}

// Creates a game and return the game ID, player ID and any errors
func (gr *GameRepo) CreateGame(gameSettings *gameLogic.GameSettings, playerName string) (uuid.UUID, uuid.UUID, error) {
	gr.lock.Lock()
	defer gr.lock.Unlock()

	game, err := gameLogic.NewGame(gameSettings, playerName)
	if err != nil {
		logger.Logger.Error("Cannot create game", "err", err)
		return uuid.UUID{}, uuid.UUID{}, err
	}

	gid := game.Id
	gr.GameMap[gid] = game
	gr.GameAgeMap[gid] = game.CreationTime

	go AddGame()
	go AddUser()

	logger.Logger.Info("Created game for", playerName)
	return gid, game.GameOwnerId, nil
}

func (gr *GameRepo) RemoveGame(gameId uuid.UUID) error {
	gr.lock.Lock()
	defer gr.lock.Unlock()

	return gr.RemoveGame(gameId)
}

// Not thread safe, to be used internally
func (gr *GameRepo) removeGame(gameId uuid.UUID) error {
	_, found := gr.GameMap[gameId]
	if !found {
		return errors.New("Cannot find game")
	}

	delete(gr.GameMap, gameId)
	delete(gr.GameAgeMap, gameId)
	return nil
}

func (gr *GameRepo) PlayerLeaveGame(gameId, playerId uuid.UUID) (gameLogic.PlayerRemovalResult, error) {
	gr.lock.Lock()
	defer gr.lock.Unlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return gameLogic.PlayerRemovalResult{}, errors.New("Cannot find game")
	}

	res, err := game.RemovePlayer(playerId)
	if err != nil {
		logger.Logger.Error("Cannot remove player from game",
			"playerId", playerId,
			"gameId", gameId,
			"err", err)
		return gameLogic.PlayerRemovalResult{}, err
	}

	if res.PlayersLeft == 0 {
		logger.Logger.Infof("Game %s has no players left, deleting it", gameId)
		gr.removeGame(gameId)
	}

	return res, nil
}

func (gr *GameRepo) DisconnectPlayer(gameId, playerId uuid.UUID) error {
	gr.lock.Lock()
	defer gr.lock.Unlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return errors.New("Cannot find game")
	}

	game.Lock.Lock()
	defer game.Lock.Unlock()

	player, found := game.PlayersMap[playerId]
	if !found {
		return errors.New("Cannot find player")
	}

	player.Connected = false
	return nil
}

func (gr *GameRepo) ConnectPlayer(gameId, playerId uuid.UUID) error {
	gr.lock.Lock()
	defer gr.lock.Unlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return errors.New("Cannot find game")
	}

	game.Lock.Lock()
	defer game.Lock.Unlock()

	player, found := game.PlayersMap[playerId]
	if !found {
		return errors.New("Cannot find player")
	}

	player.Connected = true
	return nil
}

func (gr *GameRepo) GetGames() []*gameLogic.Game {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	length := len(gr.GameMap)
	games := make([]*gameLogic.Game, length)

	i := 0
	for _, game := range gr.GameMap {
		games[i] = game
		i++
	}

	return games
}

func (gr *GameRepo) JoinGame(gameId, playerId uuid.UUID, password string) error {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		msg := fmt.Sprintf("Cannot find game with id %s", gameId)
		logger.Logger.Error(msg)
		return errors.New(msg)
	}

	if game.Settings.Password != password {
		return errors.New("Incorrect password")
	}

	_, found = game.PlayersMap[playerId]
	if !found {
		logger.Logger.Error("Cannot find player in game",
			"playerId", playerId,
			"gameId", gameId)
		msg := fmt.Sprintf("Cannot find player with id %s in game with id %s", playerId, gameId)
		return errors.New(msg)
	}
	return nil
}

func (gr *GameRepo) GetGame(gameId uuid.UUID) (*gameLogic.Game, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return nil, errors.New("Cannot find game")
	}
	return game, nil
}

func (gr *GameRepo) CreatePlayer(gameId uuid.UUID, playerName, password string) (uuid.UUID, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return uuid.UUID{}, errors.New("Cannot find game")
	}

	if game.Settings.Password != password {
		return uuid.UUID{}, errors.New("Incorrect password")
	}

	playerId, err := game.AddPlayer(playerName)
	if err != nil {
		return uuid.UUID{}, err
	}

	go AddUser()
	return playerId, nil
}

func (gr *GameRepo) GetPlayerName(gameId, playerId uuid.UUID) (string, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return "", errors.New("Cannot find game")
	}

	game.Lock.Lock()
	defer game.Lock.Unlock()

	player, found := game.PlayersMap[playerId]
	if !found {
		return "", errors.New("Cannot find player")
	}

	return player.Name, nil
}

func (gr *GameRepo) ChangeSettings(gameId uuid.UUID, settings gameLogic.GameSettings) error {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return errors.New("Cannot find game")
	}

	err := game.ChangeSettings(settings)
	if err != nil {
		return err
	}
	return nil
}
