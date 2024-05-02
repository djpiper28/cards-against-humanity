package gameLogic

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/logger"
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

	HandSize = 7
)

// Game settings used for the internal state and game creation
type GameSettings struct {
	// Game ends when this amount of rounds is reached
	MaxRounds uint `json:"maxRounds"`
	// Game ends when someone reaches this amount of points
	PlayingToPoints uint `json:"playingToPoints"`
	// Allows a game to have a password, this will be stored in plaintext like a chad
	// Empty string is no password
	Password   string      `json:"gamePassword"`
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
		logger.Logger.Errorf("Max rounds (%d) is less than min rounds (%d)",
			gs.MaxRounds,
			MinRounds)
		return false
	}

	if gs.MaxRounds > MaxRounds {
		logger.Logger.Errorf("Max Rounds (%d) is greater than max rounds (%d)",
			gs.MaxRounds,
			MaxRounds)
		return false
	}

	if gs.PlayingToPoints < MinPlayingToPoints {
		logger.Logger.Errorf("Playing to points (%d) is less than min playing to points (%d)",
			gs.PlayingToPoints,
			MinPlayingToPoints)
		return false
	}

	if gs.PlayingToPoints > MaxPlayingToPoints {
		logger.Logger.Errorf("Playing to points (%d) is greater than max playing to points (%d)",
			gs.PlayingToPoints,
			MaxPlayingToPoints)
		return false
	}

	if len(gs.Password) > MaxPasswordLength {
		logger.Logger.Errorf("Game password (%d) is too long (%d))",
			len(gs.Password),
			MaxPasswordLength)
		return false
	}

	if gs.MaxPlayers < MinPlayers {
		logger.Logger.Errorf("Max players (%d) is less than min players (%d)",
			gs.MaxPlayers,
			MinPlayers)
		return false
	}

	if gs.MaxPlayers > MaxPlayers {
		logger.Logger.Errorf("Max players (%d) is greater than max players (%d)",
			gs.MaxPlayers,
			MaxPlayers)
		return false
	}

	if len(gs.CardPacks) < MinCardPacks {
		logger.Logger.Errorf("Card packs length (%d) is less than min card packs length (%d)",
			len(gs.CardPacks),
			MinCardPacks)
		return false
	}

	return true
}

type GameState int

const (
	GameStateInLobby GameState = iota + 1
	GameStateWhiteCardsBeingSelected
	GameStateCzarJudgingCards
	GameStateDisplayingWinningCard
)

type Game struct {
	Id         uuid.UUID
	Players    []uuid.UUID
	PlayersMap map[uuid.UUID]*Player

	CurrentCardCzarId uuid.UUID
	GameOwnerId       uuid.UUID

	CurrentRound     uint
	Settings         *GameSettings
	CurrentBlackCard *BlackCard
	CardDeck         *CardDeck
	CreationTime     time.Time
	GameState        GameState
	Lock             sync.Mutex
}

func NewGame(gameSettings *GameSettings, hostPlayerName string) (*Game, error) {
	if !gameSettings.Validate() {
		return nil, errors.New("Cannot validate the game settings")
	}

	hostPlayer, err := NewPlayer(hostPlayerName)
	if err != nil {
		logger.Logger.Error("Cannot create game due to an error making the player",
			"err", err)
		return nil, err
	}

	playersMap := make(map[uuid.UUID]*Player)
	playersMap[hostPlayer.Id] = hostPlayer

	players := make([]uuid.UUID, 1)
	players[0] = hostPlayer.Id

	return &Game{Id: uuid.New(),
		PlayersMap:   playersMap,
		Players:      players,
		GameOwnerId:  hostPlayer.Id,
		Settings:     gameSettings,
		CreationTime: time.Now(),
		GameState:    GameStateInLobby}, nil
}

// Information that the client sees about a game
type GameInfo struct {
	Id          uuid.UUID `json:"id"`
	PlayerCount int       `json:"playerCount"`
	MaxPlayers  uint      `json:"maxPlayers"`
	PlayingTo   uint      `json:"playingTo"`
	HasPassword bool      `json:"hasPassword"`
}

// Info about a game to see before you join it
func (g *Game) Info() GameInfo {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	return GameInfo{Id: g.Id,
		PlayerCount: len(g.Players),
		MaxPlayers:  g.Settings.MaxPlayers,
		PlayingTo:   g.Settings.PlayingToPoints,
		HasPassword: g.Settings.Password != ""}
}

// Information about a game you can see when you join. Settings - password + players
type GameStateInfo struct {
	Id       uuid.UUID    `json:"id"`
	Settings GameSettings `json:"settings"`

	CurrentRound     uint       `json:"currentRound"`
	CreationTime     time.Time  `json:"creationTime"`
	GameState        GameState  `json:"gameState"`
	CurrentBlackCard *BlackCard `json:"currentBlackCard"`

	Players []Player `json:"players"`

	CurrentCardCzarId uuid.UUID `json:"currentCardCzarId"`
	GameOwnerId       uuid.UUID `json:"gameOwnerId"`
}

// The state of a game for player who has just joined a game
// or has become de-synced
func (g *Game) StateInfo() GameStateInfo {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	players := make([]Player, len(g.Players))
	for i, playerId := range g.Players {
		players[i] = Player{
			Id:        playerId,
			Name:      g.PlayersMap[playerId].Name,
			Points:    g.PlayersMap[playerId].Points,
			Connected: g.PlayersMap[playerId].Connected,
		}
	}

	return GameStateInfo{Id: g.Id,
		Settings:          *g.Settings,
		CurrentRound:      g.CurrentRound,
		CreationTime:      g.CreationTime,
		GameState:         g.GameState,
		CurrentBlackCard:  g.CurrentBlackCard,
		Players:           players,
		CurrentCardCzarId: g.CurrentCardCzarId,
		GameOwnerId:       g.GameOwnerId,
	}
}

func (g *Game) AddPlayer(playerName string) (uuid.UUID, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if len(g.Players) >= int(g.Settings.MaxPlayers) {
		return uuid.UUID{}, errors.New("Cannot add more than max players")
	}

	player, err := NewPlayer(playerName)
	if err != nil {
		return uuid.UUID{}, errors.New(fmt.Sprintf("Cannot create player %s", err))
	}

	for _, playerId := range g.Players {
		player, _ := g.PlayersMap[playerId]
		if player == nil {
			return uuid.UUID{}, errors.New("Cannot find the player from the map within the map")
		}

		if playerName == player.Name {
			return uuid.UUID{}, errors.New("Players cannot have the same name as each other")
		}
	}

	g.Players = append(g.Players, player.Id)
	g.PlayersMap[player.Id] = player
	return player.Id, nil
}

type PlayerRemovalResult struct {
	NewGameOwner uuid.UUID
	PlayersLeft  int
}

func (g *Game) RemovePlayer(playerToRemoveId uuid.UUID) (PlayerRemovalResult, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	_, found := g.PlayersMap[playerToRemoveId]
	if !found {
		return PlayerRemovalResult{}, errors.New("Player is not in the game")
	}

	delete(g.PlayersMap, playerToRemoveId)

	players := make([]uuid.UUID, 0)
	for _, pid := range g.Players {
		if pid != playerToRemoveId {
			players = append(players, pid)
		}
	}
	g.Players = players

	// If the game owner has left, then a random player should be assigned
	res := PlayerRemovalResult{PlayersLeft: len(g.Players)}
	playersLeft := len(g.Players)
	if playerToRemoveId == g.GameOwnerId && playersLeft > 0 {
		i := rand.Intn(playersLeft)
		g.GameOwnerId = g.Players[i]
		res.NewGameOwner = g.GameOwnerId
	}
	return res, nil
}

func (g *Game) StartGame() error {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if g.GameState != GameStateInLobby {
		return errors.New("The game is not in the lobby so cannot be started")
	}

	if len(g.Players) < MinPlayers {
		return errors.New(fmt.Sprintf("Cannot start game until the minimum amount of players %d have joined the game", MinPlayers))
	}

	deck, err := AccumalateCardPacks(g.Settings.CardPacks)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot create the game deck %s", err))
	}
	g.CardDeck = deck

	blackCard, err := g.CardDeck.GetNewBlackCard()
	if err != nil {
		// Allegedly impossible to get here
		return errors.New("Cannot get a black card")
	}

	g.CurrentBlackCard = blackCard
	g.GameState = GameStateWhiteCardsBeingSelected
	g.CurrentRound++

	for _, player := range g.PlayersMap {
		cards, err := g.CardDeck.GetNewWhiteCards(HandSize)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot create game %s", err))
		}

		cardIndexSlice := make(map[int]*WhiteCard)
		for _, card := range cards {
			cardIndexSlice[card.Id] = card
		}
		player.Hand = cardIndexSlice
	}
	return nil
}

type GameMetrics struct {
	PlayersConnected int
	Players          int
}

func (g *Game) Metrics() GameMetrics {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	metrics := GameMetrics{}

	for _, player := range g.PlayersMap {
		metrics.Players++
		if player.Connected {
			metrics.PlayersConnected++
		}
	}

	return metrics
}
