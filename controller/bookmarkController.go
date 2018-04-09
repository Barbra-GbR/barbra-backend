package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/bitphinix/barbra-backend/payloads"
	"net/http"
	"github.com/bitphinix/barbra-backend/models"
	"log"
	"gopkg.in/mgo.v2/bson"
)

type BookmarkController struct{}

func (BookmarkController) AddUserBookmark(c *gin.Context) {
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
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (BookmarkController) RemoveUserBookmark(c *gin.Context) {
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
