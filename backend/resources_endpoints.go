package main

import (
	"net/http"

	"github.com/djpiper28/cards-against-humanity/backend/gameLogic"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"time"

	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
)

type getPacksResp map[uuid.UUID]*gameLogic.CardPack

// @Summary		Gets all of the resource cards packs
// @Description	Gets all of the card packs as a map
// @Tags			resources
// @Accept			json
// @Produce		json
// @Success		200	{object}	getPacksResp
// @Router			/res/packs [get]
func getPacks(c *gin.Context) {
	c.JSON(http.StatusOK, gameLogic.AllPacks)
}

func SetupResoucesEndpoints(r *gin.Engine) {
	memoryStore := persist.NewMemoryStore(time.Hour)

	resourcesRoutes := r.Group("/res")
	{
		resourcesRoutes.GET("/packs", cache.CacheByRequestURI(memoryStore, 2*time.Minute), getPacks)
	}
}
