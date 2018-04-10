package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/Barbra-GbR/barbra-backend/models"
	"github.com/Barbra-GbR/barbra-backend/helpers"
	"github.com/Barbra-GbR/barbra-backend/auth"
)

type LoginController struct{}

func (LoginController) Auth(c *gin.Context) {
	providerId := c.Param("provider")
	userManager := auth.GetUserManager();
	state, err := auth.GenerateToken(64)

	if err != nil {
		Error(c, http.StatusInternalServerError, "unable to generate login url")
		return
	}

	url, err := userManager.GenerateLoginUrl(providerId, state)

	if err != nil {
		Error(c, http.StatusInternalServerError, "unable to generate login url")
		return
	}

	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("provider_id", providerId)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (LoginController) AuthCallback(c *gin.Context) {
	session := sessions.Default(c)
	userManager := auth.GetUserManager();
	jwt := auth.GetJWT()

	providerId := c.Param("provider")
	state := c.Query("state")

	//cross-site-forgery protection
	if providerId != session.Get("provider_id") || state != session.Get("state") {
		//TODO: Error page (unknown error)
		Error(c, http.StatusBadRequest, "Error page (unknown error)")
		return
	}

	//Delete session
	session.Options(sessions.Options{MaxAge: -1})
	session.Clear()
	session.Save()

	account, err := userManager.GetAccount(providerId, c.Query("code"))

	if err == models.ErrEmailAlreadyInUse {
		//TODO: Error page (email already in use)
		Error(c, http.StatusBadRequest, "Error page (email already in use)")
		return
	}

	if err == models.ErrInvalidPayload {
		//TODO: Error page (permissions, incomplete profile)
		Error(c, http.StatusBadRequest, "Error page (permissions, incomplete profile)")
		return
	}

	if err != nil {
		//TODO: Error page (unknown error)
		Error(c, http.StatusBadRequest, "Error page (unknown error)")
		return
	}

	token, err := jwt.GenerateToken(account.Id)

	if err != nil {
		//TODO: Error page (unknown error)
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, helpers.M{"jwt-token": token})
}
