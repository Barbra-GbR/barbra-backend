package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/bitphinix/barbra_backend/models"
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
	userInfo := new(models.UserInfo)
	err := c.BindJSON(userInfo)
	if err != nil {
		Error(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	err = user.UpdateAccountInfo(userInfo)

	if err == models.ErrEmailAlreadyInUse {
		Error(c, http.StatusConflict, "email already in use")
		return
	}

	if err != nil {
		Error(c, http.StatusUnprocessableEntity, "invalid userInfo")
		return
	}

	c.JSON(http.StatusOK, user)
}
