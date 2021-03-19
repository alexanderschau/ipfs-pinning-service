package remotePinning

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	router.GET("/pins", PinsGet)
	router.POST("/pins", PinsPost)
	router.DELETE("/pins/:requestid", PinsRequestidDelete)
	router.GET("/pins/:requestid", PinsRequestidGet)
	router.POST("/pins/:requestid", PinsRequestidPost)
}
