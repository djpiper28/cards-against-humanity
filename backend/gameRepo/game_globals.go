package gameRepo

import (
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
)

var Repo *GameRepo = initGlobals()

func gameRemovalDaemon(repo *GameRepo) {
	for {
		games := repo.EndOldGames()
		for _, game := range games {
			repo.RemoveGame(game)
		}
		time.Sleep(time.Second)
	}
}

func initGlobals() *GameRepo {
	err := gameLogic.LoadPacks()
	if err != nil {
		logger.Logger.Fatal("Cannot create the card packs", "err", err)
	}

	logger.Logger.Info("Initialising Game Repo")

	instance := New()
	go gameRemovalDaemon(instance)
	return instance
}
