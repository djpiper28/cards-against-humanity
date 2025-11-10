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
	MaxPlayingToPoints = 100

	MaxPasswordLength = 20

	MinPlayers = 2
	MaxPlayers = 50

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
	CardPacks  []uuid.UUID `json:"cardPacks"`
}

func DefaultGameSettings() *GameSettings {
	return &GameSettings{MaxRounds: MaxRounds,
		PlayingToPoints: 10,
		Password:        "",
		MaxPlayers:      10,
		CardPacks:       []uuid.UUID{DefaultCardPack().Id}}
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
	// GameStateDisplayingWinningCard
	// Not used in normal game play, used to denote unchanged or error game states
	GameStateEmpty = -1
)

func (state GameState) String() string {
	switch state {
	case GameStateInLobby:
		return "In Lobby"
	case GameStateWhiteCardsBeingSelected:
		return "White Cards Being Selected"
	case GameStateCzarJudgingCards:
		return "Czar Judging Cards"
	case GameStateEmpty:
		return "Error"
	}
	return "Unknown"
}

type Game struct {
	// Used to determine if the game should be shutdown
	LastAction time.Time

	Id         uuid.UUID
	Players    []uuid.UUID
	PlayersMap map[uuid.UUID]*Player

	CurrentCardCzarId uuid.UUID
	GameOwnerId       uuid.UUID
	PreviousWinner    PreviousWinner

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

	return &Game{LastAction: time.Now(),
		Id:           uuid.New(),
		PlayersMap:   playersMap,
		Players:      players,
		GameOwnerId:  hostPlayer.Id,
		Settings:     gameSettings,
		CreationTime: time.Now(),
		GameState:    GameStateInLobby}, nil
}

func (g *Game) updateLastAction() {
	g.LastAction = time.Now()
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

type InitialRoundInfo struct {
	CardCzarId           uuid.UUID      `json:"czarId"`
	RoundNumber          uint           `json:"roundNumber"`
	BlackCard            *BlackCard     `json:"blackCard"`
	PlayersWhoHavePlayed []uuid.UUID    `json:"playersWhoHavePlayed"`
	YourHand             []*WhiteCard   `json:"yourHand"`
	YourPlays            []int          `json:"yourPlays"`
	PreviousWinner       PreviousWinner `json:"previousWinner"`
}

// Information about a game you can see when you join. Settings - password + players
type GameStateInfo struct {
	Id           uuid.UUID        `json:"id"`
	Settings     GameSettings     `json:"settings"`
	CreationTime time.Time        `json:"creationTime"`
	GameState    GameState        `json:"gameState"`
	Players      []Player         `json:"players"`
	GameOwnerId  uuid.UUID        `json:"gameOwnerId"`
	RoundInfo    InitialRoundInfo `json:"roundInfo"`
	// Only used when there it is in the judging round
	AllPlays [][]*WhiteCard `json:"allPlays"`
}

// The state of a game for player who has just joined a game
// or has become de-synced
func (g *Game) StateInfo(pid uuid.UUID) GameStateInfo {
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

	// error state does not matter
	initialRoundInfo := InitialRoundInfo{
		CardCzarId:           g.CurrentCardCzarId,
		RoundNumber:          g.CurrentRound,
		BlackCard:            g.CurrentBlackCard,
		PlayersWhoHavePlayed: make([]uuid.UUID, 0),
		YourHand:             make([]*WhiteCard, 0),
		YourPlays:            make([]int, 0),
		PreviousWinner:       g.PreviousWinner,
	}

	allPlays := make([][]*WhiteCard, 0)
	for playerId, player := range g.PlayersMap {
		if pid == playerId {
			for _, card := range player.Hand {
				initialRoundInfo.YourHand = append(initialRoundInfo.YourHand, card)
			}

			for _, card := range player.CurrentPlay {
				initialRoundInfo.YourPlays = append(initialRoundInfo.YourPlays, card.Id)
			}
		}

		if len(player.CurrentPlay) > 0 {
			initialRoundInfo.PlayersWhoHavePlayed = append(initialRoundInfo.PlayersWhoHavePlayed, playerId)
		}

		if g.GameState == GameStateCzarJudgingCards {
			allPlays = append(allPlays, player.CurrentPlay)
		}
	}

	return GameStateInfo{
		Id:           g.Id,
		Settings:     *g.Settings,
		CreationTime: g.CreationTime,
		GameState:    g.GameState,
		Players:      players,
		GameOwnerId:  g.GameOwnerId,
		RoundInfo:    initialRoundInfo,
		AllPlays:     allPlays,
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
		return uuid.UUID{}, fmt.Errorf("Cannot create player %s", err)
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
	if g.GameState != GameStateInLobby {
		g.newWhiteCards()
	}
	g.updateLastAction()
	return player.Id, nil
}

type PlayerRemovalResult struct {
	NewGameOwner        uuid.UUID
	PlayersLeft         int
	NewGameState        GameState
	CzarJudingPhaseInfo CzarJudingPhaseInfo
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

	// If the remaining players have all played move to judging
	if g.GameState == GameStateWhiteCardsBeingSelected && g.playersHaveAllPlayed() {
		czarJudgingPhaseInfo, err := g.moveToCzarJudgingPhase()
		if err != nil {
			return PlayerRemovalResult{}, err
		}

		res.NewGameState = g.GameState
		res.CzarJudingPhaseInfo = czarJudgingPhaseInfo
	}

	g.updateLastAction()
	// TODO: If there are below the minimum amount of players move to the lobby
	// TODO: If the player was the czar
	return res, nil
}

type PreviousWinner struct {
	PlayerId   uuid.UUID    `json:"playerId"`
	BlackCard  *BlackCard   `json:"blackCard"`
	Whitecards []*WhiteCard `json:"whiteCards"`
}

// This contains everyone's hands, so just remember not to send it to all players lol
type RoundInfo struct {
	PlayerHands       map[uuid.UUID][]*WhiteCard
	PlayersPlays      map[uuid.UUID][]*WhiteCard
	CurrentBlackCard  *BlackCard
	CurrentCardCzarId uuid.UUID
	RoundNumber       uint
	PreviousWinner    PreviousWinner
}

// Not thread safe
// Error state: no round in progress
func (g *Game) roundInfo() (RoundInfo, error) {
	if g.GameState != GameStateWhiteCardsBeingSelected {
		return RoundInfo{}, errors.New("The game is not in the white card selection phase")
	}

	info := RoundInfo{
		CurrentBlackCard:  g.CurrentBlackCard,
		CurrentCardCzarId: g.CurrentCardCzarId,
		RoundNumber:       g.CurrentRound,
		PlayerHands:       make(map[uuid.UUID][]*WhiteCard),
		PlayersPlays:      make(map[uuid.UUID][]*WhiteCard),
		PreviousWinner:    g.PreviousWinner,
	}

	for _, player := range g.PlayersMap {
		handCards := make([]*WhiteCard, 0)
		for _, card := range player.Hand {
			handCards = append(handCards, card)
		}
		info.PlayerHands[player.Id] = handCards

		playedCards := make([]*WhiteCard, 0)
		for _, card := range player.CurrentPlay {
			playedCards = append(playedCards, card)
		}
		info.PlayersPlays[player.Id] = playedCards
	}
	return info, nil
}

// Error state: no round in progress
func (g *Game) RoundInfo() (RoundInfo, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	return g.roundInfo()
}

func (g *Game) newBlackCard() error {
	blackCard, err := g.CardDeck.GetNewBlackCard()
	if err != nil {
		return errors.New("Cannot get a black card")
	}

	g.CurrentBlackCard = blackCard
	return nil
}

// Not thread safe, gives each player new cards
func (g *Game) newWhiteCards() error {
	for pid, p := range g.PlayersMap {
		numberOfCardsRequired := HandSize - len(p.Hand)
		if numberOfCardsRequired < 0 {
			continue
		}

		cards, err := g.CardDeck.GetNewWhiteCards(uint(numberOfCardsRequired))
		if err != nil {
			return fmt.Errorf("Cannot create game: %s", err)
		}

		for _, card := range cards {
			g.PlayersMap[pid].Hand[card.Id] = card
		}
	}

	return nil
}

func (g *Game) newCzar() uuid.UUID {
	var newCzarId uuid.UUID
	currentCzarIndex := len(g.Players) - 2
	for i, pid := range g.Players {
		if pid == g.CurrentCardCzarId {
			currentCzarIndex = i
			break
		}
	}

	newCzarId = g.Players[(currentCzarIndex+1)%len(g.Players)]
	g.CurrentCardCzarId = newCzarId
	return newCzarId
}

func (g *Game) StartGame() (RoundInfo, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if g.GameState != GameStateInLobby {
		return RoundInfo{}, errors.New("The game is not in the lobby so cannot be started")
	}

	if len(g.Players) < MinPlayers {
		return RoundInfo{}, fmt.Errorf("Cannot start game until the minimum amount of players %d have joined the game", MinPlayers)
	}

	// Validate the decks, as before this they are updated via a lazy copy
	packs, err := GetCardPacks(g.Settings.CardPacks)
	if err != nil {
		return RoundInfo{}, fmt.Errorf("Cannot create the game deck %s", err)
	}

	deck, err := AccumalateCardPacks(packs)
	if err != nil {
		return RoundInfo{}, fmt.Errorf("Cannot create the game deck %s", err)
	}
	g.CardDeck = deck
	g.CurrentRound = 1
	g.newCzar()
	g.GameState = GameStateWhiteCardsBeingSelected

	// Reset scores
	for _, player := range g.PlayersMap {
		player.Points = 0
	}

	err = g.newWhiteCards()
	if err != nil {
		return RoundInfo{}, err
	}

	err = g.newBlackCard()
	if err != nil {
		return RoundInfo{}, err
	}

	g.updateLastAction()
	return g.roundInfo()
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

func (g *Game) ChangeSettings(newSettings GameSettings) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if g.GameState != GameStateInLobby {
		return errors.New("Settings can only be changed in the lobby")
	}

	if !newSettings.Validate() {
		return errors.New("Cannot validate the new settings")
	}

	g.Settings = &newSettings
	g.updateLastAction()
	return nil
}

type PlayerHands struct {
	Hands map[uuid.UUID][]*WhiteCard
}

func (playerHands *PlayerHands) fromGame(game *Game) {
	newHands := make(map[uuid.UUID][]*WhiteCard)
	for pid, player := range game.PlayersMap {
		hand := make([]*WhiteCard, 0)
		for _, card := range player.Hand {
			hand = append(hand, card)
		}

		newHands[pid] = hand
	}

	playerHands.Hands = newHands
}

type CzarJudingPhaseInfo struct {
	AllPlays [][]*WhiteCard
	PlayerHands
}

func (g *Game) checkState(expectedState GameState) bool {
	if g.GameState != expectedState {
		logger.Logger.Warn("Attempted to run a command whilst in the wrong state",
			"currentState", g.GameState,
			"expectedState", expectedState)
		return false
	}
	return true
}

// This assumes that all players have played, please sanity check before calling,
// see PlayCards.
// Not thread safe.
func (g *Game) moveToCzarJudgingPhase() (CzarJudingPhaseInfo, error) {
	if !g.checkState(GameStateWhiteCardsBeingSelected) {
		return CzarJudingPhaseInfo{}, errors.New("The game is not in the white card selection phase")
	}

	// Remove cards from each players hand, and add to allPlays
	g.GameState = GameStateCzarJudgingCards
	allPlays := make([][]*WhiteCard, 0)
	for _, player := range g.PlayersMap {
		if player.CurrentPlay == nil {
			continue
		}

		playersPlay := make([]*WhiteCard, 0)
		for _, card := range player.CurrentPlay {
			playersPlay = append(playersPlay, card)
		}

		if len(playersPlay) > 0 {
			allPlays = append(allPlays, playersPlay)
		}

		newHand := make(map[int]*WhiteCard)
		for _, playerCard := range player.Hand {
			// Remove cards in current play from the hand
			found := false
			for _, card := range player.CurrentPlay {
				if card.Id == playerCard.Id {
					found = true
					break
				}
			}

			if !found {
				newHand[playerCard.Id] = playerCard
			}
		}

		player.Hand = newHand
	}

	err := g.newWhiteCards()
	if err != nil {
		logger.Logger.Warn("Cannot give players new cards, after judging the game will end.")
		return CzarJudingPhaseInfo{}, nil
	}

	var playerHands PlayerHands
	playerHands.fromGame(g)

	for i := range allPlays {
		j := rand.Intn(i + 1)
		allPlays[i], allPlays[j] = allPlays[j], allPlays[i]
	}

	// Copy out players hands
	return CzarJudingPhaseInfo{
		AllPlays:    allPlays,
		PlayerHands: playerHands,
	}, nil
}

type PlayCardsResult struct {
	// Moves to the czar judging phase when all players have played
	MovedToNextCardCzarPhase bool
	CzarJudingPhaseInfo      CzarJudingPhaseInfo
}

// Whether all of the players have played.
// Not thread safe.
func (g *Game) playersHaveAllPlayed() bool {
	allPlayersPlayed := true
	for pid, player := range g.PlayersMap {
		if g.CurrentCardCzarId == pid {
			continue
		} else if len(player.CurrentPlay) == 0 {
			allPlayersPlayed = false
			break
		}
	}
	return allPlayersPlayed
}

func (g *Game) PlayCards(playerId uuid.UUID, cardIds []int) (PlayCardsResult, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if !g.checkState(GameStateWhiteCardsBeingSelected) {
		return PlayCardsResult{}, errors.New("The game is not in the white card selection phase")
	}

	if g.CurrentCardCzarId == playerId {
		return PlayCardsResult{}, errors.New("The card czar cannot play cards")
	}

	if uint(len(cardIds)) != g.CurrentBlackCard.CardsToPlay {
		return PlayCardsResult{}, errors.New("The amount of cards played does not match the amount of blanks")
	}

	player, found := g.PlayersMap[playerId]
	if !found {
		return PlayCardsResult{}, errors.New("Cannot find the player in the game")
	}

	// Check that the cards are in the hand
	currentPlay := make([]*WhiteCard, 0)
	for _, cardId := range cardIds {
		whiteCard, err := GetWhiteCard(cardId)
		if err != nil {
			return PlayCardsResult{}, fmt.Errorf("Cannot find card %d", cardId)
		}

		currentPlay = append(currentPlay, whiteCard)
	}

	// Check that: there are no duplicates and all cards are in the players hands
	checkedCards := make(map[int]bool)
	for _, cardId := range cardIds {
		// Duplicate check
		_, found := checkedCards[cardId]
		if found {
			return PlayCardsResult{}, fmt.Errorf("Duplicate card in cardIds detected %d", cardId)
		}

		// Search for card in hand
		found = false
		checkedCards[cardId] = true
		for _, handCard := range player.Hand {
			if handCard.Id == cardId {
				found = true
				break
			}
		}

		if !found {
			return PlayCardsResult{}, fmt.Errorf("Cannot find card in your hand %d", cardId)
		}
	}

	// Save their current play
	player.CurrentPlay = currentPlay

	// Check that all players have played
	allPlayersPlayed := g.playersHaveAllPlayed()

	ret := PlayCardsResult{MovedToNextCardCzarPhase: allPlayersPlayed}
	if allPlayersPlayed {
		info, err := g.moveToCzarJudgingPhase()
		if err != nil {
			return PlayCardsResult{}, err
		}

		ret.CzarJudingPhaseInfo = info
	}

	g.updateLastAction()
	return ret, nil
}

func (g *Game) TimeSinceLastAction() time.Duration {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	return time.Since(g.LastAction)
}

type CzarSelectCardResult struct {
	NewCzarId    uuid.UUID
	WinnerId     uuid.UUID
	NewBlackCard *BlackCard
	// If there are no more black cards then the game is over
	GameEnded      bool           `json:"gameEnded"`
	PreviousWinner PreviousWinner `json:"previousWinner"`
	PlayerHands
}

func IsPlayEqual(playersPlays []*WhiteCard, otherPlays []int) bool {
	if len(playersPlays) != len(otherPlays) {
		return false
	}

	for _, card := range playersPlays {
		found := false
		for _, otherCard := range otherPlays {
			if otherCard == card.Id {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (g *Game) endGame() {
	logger.Logger.Info("Ending game", "gameId", g.Id)
	g.GameState = GameStateInLobby
}

// Moves to the next round, returns if the game ended
// 1. Changes to the white card selection stage
// 2. Bumps the round counter
// 3. Allocates a new czar
// 4. Makes sure each player has 7 cards
// false means the game did not end
// true means the game did end (i.e: no more cards)
func (g *Game) nextRound() bool {
	logger.Logger.Info("Moving game to next round",
		"gameId", g.Id,
		"roundNumber (pre)", g.CurrentRound)
	g.GameState = GameStateWhiteCardsBeingSelected
	g.CurrentRound += 1
	g.newCzar()

	for _, player := range g.PlayersMap {
		player.CurrentPlay = make([]*WhiteCard, 0)
	}

	err := g.newWhiteCards()
	if err != nil {
		logger.Logger.Warnf("Cannot get a new white card for the next round %s", err)
		return true
	}

	err = g.newBlackCard()
	if err != nil {
		logger.Logger.Warnf("Cannot get a new black card for the next round %s", err)
		return true
	}

	return false
}

func (g *Game) CzarSelectCards(pid uuid.UUID, cards []int) (CzarSelectCardResult, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if !g.checkState(GameStateCzarJudgingCards) {
		return CzarSelectCardResult{}, errors.New("Not in judging phase")
	}

	if g.CurrentCardCzarId != pid {
		return CzarSelectCardResult{}, fmt.Errorf("%s is not the card czar", pid)
	}

	if len(cards) != int(g.CurrentBlackCard.CardsToPlay) {
		return CzarSelectCardResult{}, fmt.Errorf("Expected %d cards, found %d",
			g.CurrentBlackCard.CardsToPlay,
			len(cards))
	}

	// Find the player who won
	var winnerId uuid.UUID
	for _, player := range g.PlayersMap {
		if IsPlayEqual(player.CurrentPlay, cards) {
			winnerId = player.Id
			break
		}
	}

	var nullId uuid.UUID
	if winnerId == nullId {
		return CzarSelectCardResult{}, errors.New("Cannot find a player that has played those cards")
	}

	// Produce result
	g.PlayersMap[winnerId].Points += 1
	g.PreviousWinner.BlackCard = g.CurrentBlackCard
	g.PreviousWinner.PlayerId = winnerId

	whiteCards := make([]*WhiteCard, len(cards))
	for i, cardId := range cards {
		card, err := GetWhiteCard(cardId)
		if err != nil {
			logger.Logger.Error("Cannot find card", "err", err, "cardId", cardId)
			continue
		}

		whiteCards[i] = card
	}
	g.PreviousWinner.Whitecards = whiteCards

	endGame := g.nextRound()

	if g.CurrentRound >= g.Settings.MaxRounds {
		endGame = true
	}

	for _, player := range g.PlayersMap {
		if player.Points >= int(g.Settings.PlayingToPoints) {
			endGame = true
			break
		}
	}

	if endGame {
		g.endGame()
	}

	var playerHands PlayerHands
	playerHands.fromGame(g)

	g.updateLastAction()
	return CzarSelectCardResult{
		WinnerId:       winnerId,
		NewBlackCard:   g.CurrentBlackCard,
		NewCzarId:      g.CurrentCardCzarId,
		GameEnded:      endGame,
		PlayerHands:    playerHands,
		PreviousWinner: g.PreviousWinner,
	}, nil
}

func (g *Game) SkipBlackCard(pid uuid.UUID) (*BlackCard, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	if pid != g.CurrentCardCzarId {
		return nil, errors.New("Player is not the card czar so cannot skip the black card")
	}

	if !g.checkState(GameStateWhiteCardsBeingSelected) {
		return nil, errors.New("Cannot skip a card unless white cards are being selected")
	}

	err := g.newBlackCard()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot skip black card"), err)
	}

	// Clear plays
	for _, pid := range g.Players {
		player, _ := g.PlayersMap[pid]
		player.CurrentPlay = make([]*WhiteCard, 0)
	}

	logger.Logger.Info("Skipped the black card", "gameId", g.Id, "playerId", pid)
	return g.CurrentBlackCard, nil
}

func (g *Game) MulliganHand(pid uuid.UUID) ([]*WhiteCard, error) {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	player, found := g.PlayersMap[pid]
	if !found {
		return nil, errors.New("Cannot find the player")
	}

	if !g.checkState(GameStateWhiteCardsBeingSelected) {
		return nil, errors.New("Not in white card selection phase")
	}

	cards, err := g.CardDeck.GetNewWhiteCards(7)
	if err != nil {
		return nil, errors.New("Cannot get new cards")
	}

	player.Hand = make(map[int]*WhiteCard, 7)
	for _, card := range cards {
		player.Hand[card.Id] = card
	}

	player.CurrentPlay = make([]*WhiteCard, 0)

	return cards, nil
}
