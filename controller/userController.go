package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/bitphinix/barbra-backend/models"
	"github.com/bitphinix/barbra-backend/payloads"
)

type UserController struct{}

func (UserController) GetAccount(c *gin.Context) {
	account, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, account)
}

func (UserController) UpdateProfile(c *gin.Context) {
	payload := new(payloads.ProfilePayload)
	err := c.BindJSON(payload)
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	err = user.UpdateProfile(payload)

	if err == models.ErrEmailAlreadyInUse {
		Error(c, http.StatusConflict, "email already in use")
		return
	}

	if err != nil {
		Error(c, http.StatusUnprocessableEntity, "invalid payload")
		return
	}
}

