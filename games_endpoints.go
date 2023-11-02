package main

import (
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

func SetupGamesEndpoints(r *gin.Engine) {
	gamesRoute := r.Group("/games")
	{
		gamesRoute.GET("/games", getGames)
	}
}
