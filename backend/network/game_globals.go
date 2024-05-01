package network

import (
	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
)

var GameRepo *gameRepo.GameRepo = initGlobals()

func initGlobals() *gameRepo.GameRepo {
	err := gameLogic.LoadPacks()
	if err != nil {
		logger.Logger.Fatal("Cannot create the card packs", "err", err)
	}

	logger.Logger.Info("Initialising Game Repo")
	return gameRepo.New()
}
