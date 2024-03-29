package main

import (
	"log"
	"net/http"

	gpmiddleware "github.com/carousell/gin-prometheus-middleware"
	docs "github.com/djpiper28/cards-against-humanity/backend/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// A function to start for use in e2e testing
func Start() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	log.Println("Starting up Cards Against Humanity server")
	router := gin.Default()

	// Setup swagger
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// System routes
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"healthy": true,
		})
	})

	// Add metrics for the amount of games and users
	router.GET("/game-metrics", func(c *gin.Context) {
		c.TOML(http.StatusOK, GetMetrics())
	})

	// Add metrics for the performance of Gin
	prometheusEndpoint := gpmiddleware.NewPrometheus("cards_against_humanity")
	prometheusEndpoint.Use(router)

	// Add CORS middleware, this also allows for other clients to use this server
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "DELETE", "GET"},
	}))

	// Setup all endpoints
	SetupGamesEndpoints(router)
	SetupResoucesEndpoints(router)

	router.Run()
}

// @title			Cards Against Humanity API
// @version		1.0
// @description	A FOSS Cards Against Humanity server written in Go
// @contact.name	Danny Piper (djpiper28)
// @contact.url	https://github.com/djpiper28/cards-against-humanity
// @contact.email	djpiper28@gmail.com
// @license.name	GNU GPL 3
// @license.url	https://github.com/djpiper28/cards-against-humanity
// @schemes		http https
// @host			localhost:8080
// @BasePath		/
func main() {
	Start()
}
