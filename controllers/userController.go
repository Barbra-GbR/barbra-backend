package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/Barbra-GbR/barbra-backend/models"
	"github.com/Barbra-GbR/barbra-backend/payloads"
	"gopkg.in/mgo.v2/bson"
	"log"
)

//Provides an webInterface for managing user information
type UserController struct{}

//Gets a users account
func (UserController) GetAccount(c *gin.Context) {
	account, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, account)
}

//Updates a users profile
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

//Adds a bookmark to a users profile
func (controller *UserController) AddBookmark(c *gin.Context) {
	payload := new(payloads.BookmarkPayload)
	err := c.BindJSON(payload)
	if err != nil || !bson.IsObjectIdHex(payload.SuggestionId) {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	container, err := account.GetBookmarkContainer()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = container.AddBookmark(bson.ObjectIdHex(payload.SuggestionId))
	if err == models.ErrSuggestionNotFound {
		Error(c, http.StatusUnprocessableEntity, "no suggestion with id")
		return
	}
	if err == models.ErrBookmarkExists {
		c.AbortWithStatus(http.StatusAlreadyReported)
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

//Removes a bookmark from a users profile
func (controller *UserController) RemoveBookmark(c *gin.Context) {
	payload := new(payloads.BookmarkPayload)
	err := c.BindJSON(payload)
	if err != nil || !bson.IsObjectIdHex(payload.SuggestionId) {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	container, err := account.GetBookmarkContainer()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = container.RemoveBookmark(bson.ObjectIdHex(payload.SuggestionId))
	if err == models.ErrSuggestionNotFound {
		Error(c, http.StatusUnprocessableEntity, "no bookmark with id")
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
