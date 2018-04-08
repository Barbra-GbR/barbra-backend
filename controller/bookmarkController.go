package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/bitphinix/barbra_backend/payloads"
	"net/http"
	"github.com/bitphinix/barbra_backend/models"
)

type BookmarkController struct {}

func (BookmarkController) AddUserBookmark(c *gin.Context) {
	payload := new(payloads.BookmarkPayload)
	err := c.BindJSON(payload)

	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	err = user.AddBookmark(payload.SuggestionId)

	if err == models.ErrSuggestionNotFound {
		Error(c, http.StatusUnprocessableEntity, "no suggestion with id")
		return
	}

	if err != nil {
		c.AbortWithStatus(http.StatusAlreadyReported)
	}
}

func (BookmarkController) RemoveUserBookmark(c *gin.Context) {
	payload := new(payloads.BookmarkPayload)
	err := c.BindJSON(payload)

	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	err = user.RemoveBookmark(payload.SuggestionId)

	if err != nil {
		Error(c, http.StatusUnprocessableEntity, "no bookmark with id")
	}
}

func (BookmarkController) GetUserBookmarks(c *gin.Context) {
	account, err := GetCurrentAccount(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, account.Bookmarks)
}
