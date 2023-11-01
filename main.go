package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC)
	log.Println("Starting up Cards Against Humanity server")

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
