package api

import (
	"net/http"

	remotePinning "github.com/alexanderschau/ipfs-pinning-service/remote-pinning"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	remotePinning.RegisterRoutes(router)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	return router
}
