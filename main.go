package main

import (
	"log"
	"net/http"

	gpmiddleware "github.com/carousell/gin-prometheus-middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	log.Println("Starting up Cards Against Humanity server")
	InitGlobals()

	r := gin.Default()

	// System routes
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"healthy": true,
		})
	})

	// Add metrics for the amount of games and users
	r.GET("/game-metrics", func(c *gin.Context) {
		c.TOML(http.StatusOK, GetMetrics())
	})

	// Setup all endpoints
	SetupGamesEndpoints(r)

	// Add metrics for the performance of Gin
	p := gpmiddleware.NewPrometheus("gin")
	p.Use(r)

	r.Run()
}
