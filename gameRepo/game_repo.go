package gameRepo

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/google/uuid"
)

const (
	MaxGameInProgressAge    = time.Hour * 3
	MaxGameInLobbyAge       = time.Minute * 15
	MaxGameWithNoPlayersAge = time.Second * 2
)

type GameListPtr *gameLogic.Game

type GameRepo struct {
	// A sorted list of games, where the first game is the oldest, when a game starts it is sent to the end
	// When a game ends it is put back to the front. This is used for O(k) lookup of games to kill
	GamesByAge *list.List
	GameMap    map[uuid.UUID]*gameLogic.Game
	GameAgeMap map[uuid.UUID]time.Time
	lock       sync.RWMutex
}

func New() *GameRepo {
	return &GameRepo{GamesByAge: list.New(), GameMap: make(map[uuid.UUID]*gameLogic.Game), GameAgeMap: make(map[uuid.UUID]time.Time)}
}

// Creates a game and return the game ID, player ID and any errors
func (gr *GameRepo) CreateGame(gameSettings *gameLogic.GameSettings, playerName string) (uuid.UUID, uuid.UUID, error) {
	gr.lock.Lock()
	defer gr.lock.Unlock()

	game, err := gameLogic.NewGame(gameSettings, playerName)
	if err != nil {
		log.Println("Cannot create game", err)
		return uuid.UUID{}, uuid.UUID{}, err
	}

	id := game.Id
	gr.GamesByAge.PushBack(GameListPtr(game))
	gr.GameMap[id] = game
	gr.GameAgeMap[id] = game.CreationTime

	log.Println("Created game for", playerName)
	return id, game.GameOwnerId, nil
}

func (gr *GameRepo) GetGames() []*gameLogic.Game {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	length := gr.GamesByAge.Len()
	games := make([]*gameLogic.Game, length)

	current := gr.GamesByAge.Front()
	for i := 0; i < length; i++ {
		games[i] = current.Value.(GameListPtr)
		current = current.Next()
	}

	return games
}

func (gr *GameRepo) JoinGame(gameId, playerId uuid.UUID) error {
	gr.lock.RLock()
	defer gr.lock.RUnlock()

	game, found := gr.GameMap[gameId]
	if !found {
		msg := fmt.Sprintf("Cannot find game with id %s", gameId)
		log.Println(msg)
		return errors.New(msg)
	}

	_, found = game.PlayersMap[playerId]
	if !found {
		msg := fmt.Sprintf("Cannot find player with id %s in game with id %s", playerId, gameId)
		log.Println(msg)
		return errors.New(msg)
	}
	return nil
}
