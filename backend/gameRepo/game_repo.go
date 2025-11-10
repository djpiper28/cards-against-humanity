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
	MaxGameAge = time.Minute * 20
	// Allow for enough time for players to reconnect in the case of a network drop
	MaxGameWithNoPlayersAge = time.Minute * 2
)

type GameRepo struct {
	GameMap map[uuid.UUID]*gameLogic.Game
	lock    sync.RWMutex
}

func New() *GameRepo {
	return &GameRepo{GameMap: make(map[uuid.UUID]*gameLogic.Game)}
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

	go AddGame()
	go AddUser()

	logger.Logger.Info("Created game for", "player", playerName)
	return gid, game.GameOwnerId, nil
}

func (gr *GameRepo) RemoveGame(gameId uuid.UUID) error {
	gr.lock.Lock()
	defer gr.lock.Unlock()

	return gr.removeGame(gameId)
}

// Not thread safe, to be used internally
func (gr *GameRepo) removeGame(gameId uuid.UUID) error {
	_, found := gr.GameMap[gameId]
	if !found {
		return errors.New("Cannot find game")
	}

	delete(gr.GameMap, gameId)
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
		logger.Logger.Info("Game has no players left, deleting it",
			"gameId", gameId)
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

func (gr *GameRepo) StartGame(gameId uuid.UUID) (gameLogic.RoundInfo, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return gameLogic.RoundInfo{}, errors.New("Cannot find game")
	}

	info, err := game.StartGame()
	if err != nil {
		return gameLogic.RoundInfo{}, err
	}

	return info, nil
}

func (gr *GameRepo) PlayerPlayCards(gameId, playerId uuid.UUID, cardIds []int) (gameLogic.PlayCardsResult, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return gameLogic.PlayCardsResult{}, errors.New("Cannot find game")
	}

	info, err := game.PlayCards(playerId, cardIds)
	if err != nil {
		return gameLogic.PlayCardsResult{}, err
	}

	return info, nil
}

func (gr *GameRepo) EndOldGames() []uuid.UUID {
	start := time.Now()

	endedGames := make([]uuid.UUID, 0)
	games := gr.GetGames()
	for _, game := range games {
		remove := false
		lastActionTime := game.TimeSinceLastAction()
		if lastActionTime > MaxGameAge {
			remove = true
		} else if game.Metrics().PlayersConnected == 0 && lastActionTime > MaxGameWithNoPlayersAge {
			remove = true
		}

		if remove {
			endedGames = append(endedGames, game.Id)
		}
	}

	count := len(endedGames)
	go AddGamePurgeData(int(time.Since(start).Microseconds()), count)

	if count > 0 {
		logger.Logger.Info("Ending old games", "count", count)
	}

	gr.lock.Lock()
	defer gr.lock.Unlock()
	for _, removedGameId := range endedGames {
		gr.removeGame(removedGameId)
	}
	return endedGames
}

func (gr *GameRepo) CzarSelectsCard(gameId, pid uuid.UUID, cards []int) (gameLogic.CzarSelectCardResult, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return gameLogic.CzarSelectCardResult{}, errors.New("Cannot find game")
	}

	return game.CzarSelectCards(pid, cards)
}

func (gr *GameRepo) CzarSkipsCard(gameId, pid uuid.UUID) (*gameLogic.BlackCard, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return nil, errors.New("Cannot find game")
	}

	return game.SkipBlackCard(pid)
}

func (gr *GameRepo) MulliganHand(gameId, pid uuid.UUID) ([]*gameLogic.WhiteCard, error) {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		return nil, errors.New("Cannot find game")
	}

	return game.MulliganHand(pid)
}
