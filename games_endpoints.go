package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/gin-gonic/gin"
)

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
  PlayerName string `json:"playerName"`
  Settings gameLogic.GameSettings `json:"settings"`
}

func createGame(c *gin.Context) {
  settingsStr, err := io.ReadAll(c.Request.Body)
  if err != nil {
    c.Error(err)
  }

  var settings gameCreateSettings
  err = json.Unmarshal(settingsStr, settings)
  if err != nil {
    c.Error(err)
  }

  gameId, playerId, err := GameRepo.CreateGame(&settings.Settings, settings.PlayerName)
  if err != nil {
    c.Error(err)
  }

  _, err = WsUpgrade(c.Writer, c.Request, gameId, playerId)
  if err != nil {
    c.Error(err)
  }
  // TODO: use the connection
}

func SetupGamesEndpoints(r *gin.Engine) {
	gamesRoute := r.Group("/games")
	{
		gamesRoute.GET("/notFull", getGames)
    gamesRoute.POST("/create", createGame)
	}
}
