package middlewares

import (
	"github.com/gin-gonic/gin"
	"../auth"
	"net/http"
)

func AuthorizationMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if len(tokenString) < 1 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwt := auth.GetJWT()
	_, err := jwt.GetUserId(tokenString)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
