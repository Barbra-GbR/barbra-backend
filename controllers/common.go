package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/Barbra-GbR/barbra-backend/helpers"
	"github.com/Barbra-GbR/barbra-backend/models"
	"errors"
	"net/http"
)

var (
	ErrContextNotSet = errors.New("context not set")
)

//Aborts the current context and sends the error serialised to the client
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, helpers.M{"error": message})
	c.Abort()
}

//Returns the user set from the authenticationMiddleware
func GetCurrentAccount(c *gin.Context) (*models.UserAccount, error) {
	accountInterface, ok := c.Get("user_account")
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, ErrContextNotSet
	}

	account, ok := accountInterface.(*models.UserAccount)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, ErrContextNotSet
	}

	return account, nil
}
