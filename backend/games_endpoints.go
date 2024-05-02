package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/djpiper28/cards-against-humanity/backend/security"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	JoinGameGameIdParam   = "gameId"
	JoinGamePlayerIdParam = "playerId"
	PasswordParam         = "password"
	AuthorizationCookie   = "authentication"
)

// @Summary		Gets all of the games that are not full
// @Description	Returns a list of the games
// @Tags			games
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]gameLogic.GameInfo
// @Router			/games/notFull [get]
func getGames(c *gin.Context) {
	games := gameRepo.Repo.GetGames()
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
	GameId              uuid.UUID `json:"gameId"`
	PlayerId            uuid.UUID `json:"playerId"`
	AuthorisationCookie string    `json:"authentication"`
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

	gameId, playerId, err := gameRepo.Repo.CreateGame(&settings, createReq.PlayerName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, NewApiError(err))
		return
	}

	token, err := security.NewToken(gameId, playerId)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, NewApiError(errors.New("Cannot create authorization token")))
		return
	}

	c.SetCookie(AuthorizationCookie, token, -1, "/", "", true, true)
	resp := GameCreatedResp{GameId: gameId, PlayerId: playerId, AuthorisationCookie: token}
	c.JSON(http.StatusCreated, resp)
}

type CreatePlayerRequest struct {
	PlayerName string    `json:"playerName"`
	GameId     uuid.UUID `json:"gameId"`
	Password   string    `json:"password"`
}

type CreatePlayerResponse struct {
	PlayerId            uuid.UUID `json:"playerId"`
	AuthorisationCookie string    `json:"authentication"`
}

// @Summary		Creates a player to allow you to join a game (first step of game joining, followed by /join ing)
// @Description	Validates the player information, then tries to add them to a game and returns their ID.
// @Tags			games
// @Accept			json
// @Produce		json
// @Param			request	body					CreatePlayerRequest	true	"Player information"
// @Success		204		{CreatePlayerResponse}	uuid				"Player ID and auth token"
// @Failure		500		{object}				ApiError
// @Failure		400		{object}				ApiError
// @Router			/games/join [post]
func createPlayerForJoining(c *gin.Context) {
	req, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Logger.Error("Cannot read body", "err", err)
		c.JSON(http.StatusInternalServerError, NewApiError(errors.New("Failed to read request body")))
		return
	}

	var createReq CreatePlayerRequest
	err = json.Unmarshal(req, &createReq)
	if err != nil {
		logger.Logger.Error("Cannot unmarshal request", "err", err)
		c.JSON(http.StatusBadRequest, NewApiError(errors.New("Invalid request")))
		return
	}

	playerId, err := gameRepo.Repo.CreatePlayer(createReq.GameId, createReq.PlayerName, createReq.Password)
	if err != nil {
		logger.Logger.Error("Cannot create player", "err", err)
		c.JSON(http.StatusInternalServerError, NewApiError(err))
		return
	}

	onCreateMessage := network.RpcOnPlayerCreateMsg{
		Id:   playerId,
		Name: createReq.PlayerName,
	}

	token, err := security.NewToken(createReq.GameId, playerId)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, NewApiError(errors.New("Cannot create authorization token")))
		return
	}

	c.SetCookie(AuthorizationCookie, token, -1, "/", "", true, true)
	c.JSON(http.StatusCreated, CreatePlayerResponse{PlayerId: playerId, AuthorisationCookie: token})

	// Transmit to the other players
	msg, err := network.EncodeRpcMessage(onCreateMessage)
	if err != nil {
		logger.Logger.Error("Error encoding message", "err", err)
	}
	go network.GlobalConnectionManager.Broadcast(createReq.GameId, msg)
}

type authenticationData struct {
	GameId   uuid.UUID
	PlayerId uuid.UUID
	Token    string
	Password string
}

// Tries to read the authentication data from the cookies of the request,
// it will also verify the token, the password is not checked
func attemptAuthentication(c *gin.Context) (authenticationData, error) {
	// Validate input
	rawGameId, err := c.Cookie(JoinGameGameIdParam)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot find cookie %s", JoinGameGameIdParam)
		logger.Logger.Error(errMsg)
		c.JSON(http.StatusBadRequest, NewApiError(errors.New(errMsg)))
		return authenticationData{}, err
	}

	gameId, err := uuid.Parse(rawGameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError(err))
		return authenticationData{}, err
	}

	rawPlayerId, err := c.Cookie(JoinGamePlayerIdParam)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot find cookie %s", JoinGamePlayerIdParam)
		logger.Logger.Error(errMsg)
		c.JSON(http.StatusBadRequest, NewApiError(errors.New(errMsg)))
		return authenticationData{}, err
	}

	playerId, err := uuid.Parse(rawPlayerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError(err))
		return authenticationData{}, err
	}

	token, err := c.Cookie(AuthorizationCookie)
	if err != nil {
		c.JSON(http.StatusForbidden, NewApiError(errors.New("No authorization token provided")))
		return authenticationData{}, err
	}

	err = security.CheckToken(gameId, playerId, token)
	if err != nil {
		logger.Logger.Errorf("Player %s in game %s tried to to join a game with invalid authorisation: %s",
			playerId,
			gameId,
			err)
		c.JSON(http.StatusForbidden, NewApiError(errors.New("Not authorized")))
		return authenticationData{}, err
	}

	password, err := c.Cookie(PasswordParam)
	if err != nil {
		logger.Logger.Error("There is no password provided")
	}

	return authenticationData{GameId: gameId, PlayerId: playerId, Token: token, Password: password}, nil
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
	authData, err := attemptAuthentication(c)
	if err != nil {
		return
	}

	// Join the game
	err = gameRepo.Repo.JoinGame(authData.GameId,
		authData.PlayerId,
		authData.Password)
	if err != nil {
		logger.Logger.Error("Cannot join game", "gameId", authData.GameId, "err", err)
		c.JSON(http.StatusNotFound, NewApiError(err))
		return
	}

	// Attempt to upgrade the websocket
	logger.Logger.Info("Upgrading connection",
		"gameId",
		authData.GameId,
		"playerId",
		authData.PlayerId)
	network.WsUpgrade(c.Writer, c.Request, authData.GameId, authData.PlayerId, network.GlobalConnectionManager)
}

// @Summary		Allows the player to leave a game
// @Description	Leaves the current game that a player is in
// @Tags			games
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500	{object}	ApiError
// @Failure		404	{object}	ApiError
// @Failure		400	{object}	ApiError
// @Router			/games/leave [delete]
func leaveGame(c *gin.Context) {
	authData, err := attemptAuthentication(c)
	if err != nil {
		return
	}

	err = network.GlobalConnectionManager.RemovePlayer(authData.GameId,
		authData.PlayerId)
	if err != nil {
		logger.Logger.Errorf("Player %s in game %s was unable to leave the game: %s",
			authData.PlayerId,
			authData.GameId,
			err)
		c.JSON(http.StatusInternalServerError, NewApiError(errors.New("Cannot leave the game")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func SetupGamesEndpoints(r *gin.Engine) {
	gamesRoute := r.Group("/games")
	{
		gamesRoute.GET("/notFull", getGames)
		gamesRoute.POST("/create", createGame)
		gamesRoute.GET("/join", joinGame)
		gamesRoute.POST("/join", createPlayerForJoining)
		gamesRoute.DELETE("/leave", leaveGame)
	}
}
