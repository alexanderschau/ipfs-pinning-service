package remotePinning

import (
	"net/http"
	"strings"

	"github.com/alexanderschau/ipfs-pinning-service/auth"
	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) (bool, string) {
	authArray := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authArray) < 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"reason":  "UNAUTHORIZED",
				"details": "Access token is missing or invalid",
			},
		})
		return false, ""
	}
	accessToken := authArray[1]
	check, user := auth.CheckAuth(accessToken)
	if !check {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"reason":  "UNAUTHORIZED",
				"details": "Access token is missing or invalid",
			},
		})
		return false, ""
	}

	return true, user
}
