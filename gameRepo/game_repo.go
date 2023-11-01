package gameRepo

import (
	"container/list"
	"log"
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
}

func New() *GameRepo {
	return &GameRepo{GamesByAge: list.New(), GameMap: make(map[uuid.UUID]*gameLogic.Game), GameAgeMap: make(map[uuid.UUID]time.Time)}
}

func (gr *GameRepo) CreateGame(gameSettings *gameLogic.GameSettings, playerName string) (uuid.UUID, error) {
	log.Println("Creating game for", playerName)
	game, err := gameLogic.NewGame(gameSettings, playerName)
	if err != nil {
		log.Println("Cannot create game", err)
		return uuid.UUID{}, err
	}

	id := game.Id
	gr.GamesByAge.PushBack(GameListPtr(game))
	gr.GameMap[id] = game
	gr.GameAgeMap[id] = game.CreationTime

	return id, nil
}
