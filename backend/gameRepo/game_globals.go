package gameRepo

import (
	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/logger"
)

var Repo *GameRepo = initGlobals()

func initGlobals() *GameRepo {
	err := gameLogic.LoadPacks()
	if err != nil {
		logger.Logger.Fatal("Cannot create the card packs", "err", err)
	}

	logger.Logger.Info("Initialising Game Repo")
	return New()
}
