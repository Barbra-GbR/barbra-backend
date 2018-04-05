package controller

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"net/http"
	"../auth"
	"../models"
)

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, helpers.M{"error": message})
	c.Abort()
}

func GetCurrentAccount(c *gin.Context) (*models.UserAccount, error) {
	tokenString := c.GetHeader("Authorization")

	jwt := auth.GetJWT()
	accountId, err := jwt.GetUserId(tokenString)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, err
	}

	user, err := models.GetUserAccount(accountId)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, err
	}

	return user, err
}
