package network

import (
	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/djpiper28/cards-against-humanity/backend/gameRepo"
	"log"
)

var GameRepo *gameRepo.GameRepo = initGlobals()

func initGlobals() *gameRepo.GameRepo {
	err := gameLogic.LoadPacks()
	if err != nil {
		log.Fatal("Cannot create the card packs", err)
	}

	log.Println("Initialising Game Repo")
	return gameRepo.New()
}
