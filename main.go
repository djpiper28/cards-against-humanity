package main

import (
	"log"
	"net/http"

	"github.com/djpiper28/cards-against-humanity/gameLogic"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC)
	log.Println("Starting up Cards Against Humanity server")
	err := gameLogic.LoadPacks()
	if err != nil {
		log.Fatal("Cannot create the card packs", err)
	}

	r := gin.Default()

	// System routes
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"healthy": true,
		})
	})

	r.GET("/metrics", func(c *gin.Context) {
		c.TOML(http.StatusOK, GetMetrics())
	})

	r.Run()
}
