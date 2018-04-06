package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/bitphinix/barbra_backend/helpers"
	"github.com/bitphinix/barbra_backend/models"
)

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, helpers.M{"error": message})
	c.Abort()
}

func GetCurrentAccount(c *gin.Context) (*models.UserAccount, error) {
	accountId := c.GetString("user_id")
	user, err := models.GetUserAccount(accountId)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, err
	}

	return user, nil
}
