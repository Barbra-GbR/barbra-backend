package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/bitphinix/barbra_backend/auth"
	"github.com/bitphinix/barbra_backend/models"
	"github.com/bitphinix/barbra_backend/controller"
)

func AuthorizationMiddleware(enrolledOnly bool) func(c*gin.Context) {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if len(tokenString) < 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		jwt := auth.GetJWT()
		accountId, err := jwt.GetUserId(tokenString)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		account, err := models.GetUserAccount(accountId)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if enrolledOnly && !account.Enrolled {
			controller.Error(c, http.StatusUnauthorized, "account isn´t enrolled")
			c.Abort()
			return
		}

		c.Set("user_account", account)
	}
}

