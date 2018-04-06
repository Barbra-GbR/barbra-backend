package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/bitphinix/babra_backend/auth"
)

func AuthorizationMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if len(tokenString) < 1 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwt := auth.GetJWT()
	userId, err := jwt.GetUserId(tokenString)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("user_id", userId)
}
