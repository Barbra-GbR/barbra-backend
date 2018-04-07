package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/bitphinix/barbra_backend/models"
	"github.com/bitphinix/barbra_backend/payloads"
)

type UserController struct{}

func (UserController) GetAccount(c *gin.Context) {
	user, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, user)
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

	err = user.UpdateAccountInfo(payload)

	if err == models.ErrEmailAlreadyInUse {
		Error(c, http.StatusConflict, "email already in use")
		return
	}

	if err != nil {
		Error(c, http.StatusUnprocessableEntity, "invalid payload")
		return
	}

	c.JSON(http.StatusOK, user)
}
