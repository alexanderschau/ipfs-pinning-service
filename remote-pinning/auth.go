package remotePinning

import (
	"net/http"
	"strings"

	"github.com/alexanderschau/ipfs-pinning-service/auth"
	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) (bool, string) {
	accessToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
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
