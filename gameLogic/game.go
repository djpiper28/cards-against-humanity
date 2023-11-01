package gameLogic

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// Idk why you would only play one round, but okay
	MinRounds = 1
	// Bro going to go mad after 100 rounds
	MaxRounds = 100

	// Shit game at 1
	MinPlayingToPoints = 2
	MaxPlayingToPoints = 50

	MaxPasswordLength = 50

	MinPlayers = 3
	MaxPlayers = 20

	MinCardPacks = 1
)

// Game settings used for the internal state and game creation
type GameSettings struct {
	// Game ends when this amount of rounds is reached
	MaxRounds uint `json:"maxRounds"`
	// Game ends when someone reaches this amount of points
	PlayingToPoints uint `json:"playingToPoints"`
	// Allows a game to have a password, this will be stored in plaintext like a chad
	// Empty string is no password
	Password   string      `json:"GamePassword"`
	MaxPlayers uint        `json:"maxPlayers"`
	CardPacks  []*CardPack `json:"cardPacks"`
}

func DefaultGameSettings() *GameSettings {
	return &GameSettings{MaxRounds: MaxRounds,
		PlayingToPoints: 10,
		Password:        "",
		MaxPlayers:      10,
		CardPacks:       []*CardPack{DefaultCardPack()}}
}

func (gs *GameSettings) Validate() bool {
	if gs.MaxRounds < MinRounds {
		return false
	}

	if gs.MaxRounds > MaxRounds {
		return false
	}

	if gs.PlayingToPoints < MinPlayingToPoints {
		return false
	}

	if gs.PlayingToPoints > MaxPlayingToPoints {
		return false
	}

	if len(gs.Password) > MaxPasswordLength {
		return false
	}

	if gs.MaxPlayers < MinPlayers {
		return false
	}

	if gs.MaxPlayers > MaxPlayers {
		return false
	}

	if len(gs.CardPacks) < MinCardPacks {
		return false
	}

	return true
}

type Game struct {
	Players           map[uuid.UUID]*Player
	CurrentCardCzarId uuid.UUID
	GameOwnerId       uuid.UUID
	CurrentRound      uint
	Settings          *GameSettings
	CreationTime      time.Time
	lock              sync.Mutex
}

func NewGame(gameSettings *GameSettings, hostPlayerName string) (*Game, error) {
	if !gameSettings.Validate() {
		return nil, errors.New("Cannot validate the game settings")
	}

	hostPlayer, err := NewPlayer(hostPlayerName)
	if err != nil {
		log.Println("Cannot create game due to an error making the player", err)
		return nil, err
	}

	players := make(map[uuid.UUID]*Player)
	players[hostPlayer.Id] = hostPlayer

	return &Game{Players: players, GameOwnerId: hostPlayer.Id, Settings: gameSettings, CreationTime: time.Now()}, nil
}
