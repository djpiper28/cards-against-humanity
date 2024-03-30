package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	JoinGameGameIdParam   = "gameId"
	JoinGamePlayerIdParam = "playerId"
	PasswordParam         = "password"
)

// @Summary		Gets all of the games that are not full
// @Description	Returns a list of the games
// @Tags			games
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]gameLogic.GameInfo
// @Router			/games/notFull [get]
func getGames(c *gin.Context) {
	games := network.GameRepo.GetGames()
	info := make([]gameLogic.GameInfo, 0, len(games))

	for _, game := range games {
		gameInfo := game.Info()
		if uint(gameInfo.PlayerCount) != gameInfo.MaxPlayers {
			info = append(info, gameInfo)
		}
	}

	c.JSON(http.StatusOK, info)
}

type GameCreateSettings struct {
	MaxRounds       uint        `json:"maxRounds"`
	PlayingToPoints uint        `json:"playingToPoints"`
	Password        string      `json:"gamePassword"`
	MaxPlayers      uint        `json:"maxPlayers"`
	CardPacks       []uuid.UUID `json:"cardPacks"`
}

type GameCreateRequest struct {
	PlayerName string             `json:"playerName"`
	Settings   GameCreateSettings `json:"settings"`
}

type GameCreatedResp struct {
	GameId   uuid.UUID `json:"gameId"`
	PlayerId uuid.UUID `json:"playerId"`
}

// @Summary		Creates a new game
// @Description	Creates a new game for you to connect to via websocket upgrade afterwards by calling /games/join.
// @Tags			games
// @Accept			json
// @Produce		json
// @Param			request	body		GameCreateRequest	true	"create settings"
// @Success		204		{object}	GameCreatedResp
// @Failure		500		{object}	ApiError
// @Router			/games/create [post]
func createGame(c *gin.Context) {
	settingsStr, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, NewApiError(err))
		return
	}

	var createReq GameCreateRequest
	err = json.Unmarshal(settingsStr, &createReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, NewApiError(err))
		return
	}

	packs, err := gameLogic.GetCardPacks(createReq.Settings.CardPacks)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, NewApiError(err))
		return
	}

	settings := gameLogic.GameSettings{
		MaxRounds:       createReq.Settings.MaxRounds,
		MaxPlayers:      createReq.Settings.MaxPlayers,
		PlayingToPoints: createReq.Settings.PlayingToPoints,
		Password:        createReq.Settings.Password,
		CardPacks:       packs,
	}

	gameId, playerId, err := network.GameRepo.CreateGame(&settings, createReq.PlayerName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, NewApiError(err))
		return
	}

	resp := GameCreatedResp{GameId: gameId, PlayerId: playerId}
	c.JSON(http.StatusCreated, resp)
}

type CreatePlayerRequest struct {
	PlayerName string    `json:"playerName"`
	GameId     uuid.UUID `json:"gameId"`
}

// @Summary		Creates a player to allow you to join a game (first step of game joining, followed by /join ing)
// @Description	Validates the player information, then tries to add them to a game and returns their ID.
// @Tags			games
// @Accept			json
// @Produce		json
// @Param			request	body		CreatePlayerRequest	true	"Player information"
// @Success		200		{string}	uuid				"Player ID"
// @Failure		500		{object}	ApiError
// @Failure		400		{object}	ApiError
// @Router			/games/join [post]
func createPlayerForJoining(c *gin.Context) {
	req, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Cannot read body: %s", err)
		c.JSON(http.StatusInternalServerError, NewApiError(errors.New("Failed to read request body")))
		return
	}

	var createReq CreatePlayerRequest
	err = json.Unmarshal(req, &createReq)
	if err != nil {
		log.Printf("Cannot unmarshal request: %s", err)
		c.JSON(http.StatusBadRequest, NewApiError(errors.New("Invalid request")))
		return
	}

	password, err := c.Cookie(PasswordParam)
	if err != nil {
		log.Print("There is no password provided")
	}

	playerId, err := network.GameRepo.CreatePlayer(createReq.GameId, createReq.PlayerName, password)
	if err != nil {
		log.Printf("Cannot create player: %s", err)
		c.JSON(http.StatusInternalServerError, NewApiError(err))
		return
	}

	c.JSON(http.StatusCreated, playerId)
}

// @Summary		Joins a game and upgrades the connection to a websocket if all is well
// @Description	Validates the input, checks the game exists then tries to upgrade the socket and register the connection. See the RPC docs for what to expect on the websocket. Use playerId, password and gameId cookies to authenticate.
// @Tags			games
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500	{object}	ApiError
// @Failure		404	{object}	ApiError
// @Failure		400	{object}	ApiError
// @Router			/games/join [get]
func joinGame(c *gin.Context) {
	// Validate input
	rawGameId, err := c.Cookie(JoinGameGameIdParam)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot find cookie %s", JoinGameGameIdParam)
		log.Print(errMsg)
		c.JSON(http.StatusBadRequest, NewApiError(errors.New(errMsg)))
		return
	}

	gameId, err := uuid.Parse(rawGameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError(err))
		return
	}

	rawPlayerId, err := c.Cookie(JoinGamePlayerIdParam)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot find cookie %s", JoinGamePlayerIdParam)
		log.Print(errMsg)
		c.JSON(http.StatusBadRequest, NewApiError(errors.New(errMsg)))
		return
	}

	playerId, err := uuid.Parse(rawPlayerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError(err))
		return
	}

	password, err := c.Cookie(PasswordParam)
	if err != nil {
		log.Print("There is no password provided")
	}

	// Join the game
	err = network.GameRepo.JoinGame(gameId, playerId, password)
	if err != nil {
		log.Printf("Cannot join game (%s): %s", gameId, err)
		c.JSON(http.StatusNotFound, NewApiError(err))
		return
	}

	// Attempt to upgrade the websocket
  network.WsUpgrade(c.Writer, c.Request, gameId, playerId, network.GlobalConnectionManager)
}

func SetupGamesEndpoints(r *gin.Engine) {
	gamesRoute := r.Group("/games")
	{
		gamesRoute.GET("/notFull", getGames)
		gamesRoute.POST("/create", createGame)
		gamesRoute.GET("/join", joinGame)
		gamesRoute.POST("/join", createPlayerForJoining)
	}
}
