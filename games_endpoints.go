package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Gets all of the games that are not full
// @Description	Returns a list of the games
// @Tags			games
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]gameLogic.GameInfo
// @Router			/games/notFull [get]
func getGames(c *gin.Context) {
	games := GameRepo.GetGames()
	info := make([]gameLogic.GameInfo, 0, len(games))

	for _, game := range games {
		gameInfo := game.Info()
		if uint(gameInfo.PlayerCount) != gameInfo.MaxPlayers {
			info = append(info, gameInfo)
		}
	}

	c.JSON(http.StatusOK, info)
}

type gameCreateSettings struct {
	PlayerName string                 `json:"playerName"`
	Settings   gameLogic.GameSettings `json:"settings"`
}

type gameCreatedResp struct {
	GameId   uuid.UUID `json:"gameId"`
	PlayerId uuid.UUID `json:"playerId"`
}

// @Summary		Creates a new game
// @Description	Creates a new game for you to connect to via websocket afterwards
// @Tags			games
// @Accept			json
// @Produce		json
// @Param			request	body		gameCreateSettings	true	"create settings"
// @Success		204		{object}	gameCreatedResp
// @Router			/games/create [post]
func createGame(c *gin.Context) {
	settingsStr, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
	}

	var settings gameCreateSettings
	err = json.Unmarshal(settingsStr, &settings)
	if err != nil {
		c.Error(err)
	}

	gameId, playerId, err := GameRepo.CreateGame(&settings.Settings, settings.PlayerName)
	if err != nil {
		c.Error(err)
	}

	resp := gameCreatedResp{GameId: gameId, PlayerId: playerId}
	c.JSON(http.StatusCreated, resp)
}

func SetupGamesEndpoints(r *gin.Engine) {
	gamesRoute := r.Group("/games")
	{
		gamesRoute.GET("/notFull", getGames)
		gamesRoute.POST("/create", createGame)
	}
}
