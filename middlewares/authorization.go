package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/Barbra-GbR/barbra-backend/auth"
	"github.com/Barbra-GbR/barbra-backend/models"
	"github.com/Barbra-GbR/barbra-backend/controllers"
	"gopkg.in/mgo.v2/bson"
)

//Checks if the Authorization is valid and adds the corresponding userAccount the the config as user_account
//If the enrolledOnly flag is set only enrolled users will pass the authentication
func AuthorizationMiddleware(enrolledOnly bool) func(c*gin.Context) {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) < 1 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		jwt := auth.GetJWT()
		accountId, err := jwt.GetAccountId(tokenString)
		if err != nil || !bson.IsObjectIdHex(accountId) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		account, err := models.GetUserAccountById(bson.ObjectIdHex(accountId))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if enrolledOnly && !account.Enrolled {
			controllers.Error(c, http.StatusUnauthorized, "account isn´t enrolled")
			c.Abort()
			return
		}

		c.Set("user_account", account)
	}
}

