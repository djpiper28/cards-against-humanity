package main

import (
	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/djpiper28/cards-against-humanity/gameRepo"
	"log"
)

var GameRepo *gameRepo.GameRepo

func InitGlobals() {
	err := gameLogic.LoadPacks()
	if err != nil {
		log.Fatal("Cannot create the card packs", err)
	}

	log.Println("Initialising Game Repo")
	GameRepo = gameRepo.New()
}
