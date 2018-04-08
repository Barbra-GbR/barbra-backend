package models

import (
	"time"
	"errors"
)

var (
	ErrBookmarkExsists  = errors.New("bookmark: already exists")
	ErrBookmarkNotFound = errors.New("bookmark: bookmark not found")
)

type Bookmark struct {
	Creation     time.Time `json:"creation"      bson:"creation"                             binding:"required"`
	SuggestionId string    `json:"suggestion_id" bson:"suggestion_id" validate:"hexadecimal" binding:"required"`
}

func NewBookmark(suggestionId string) *Bookmark {
	return &Bookmark{
		SuggestionId: suggestionId,
		Creation:     time.Now(),
	}
}

func (account *UserAccount) AddBookmark(suggestionId string) error {
	if _, ok := account.Bookmarks[suggestionId]; ok {
		return ErrBookmarkExsists
	}

	if !SuggestionExists(suggestionId) {
		return ErrSuggestionNotFound
	}

	bookmark := NewBookmark(suggestionId)
	account.Bookmarks[suggestionId] = bookmark
	account.Update()

	return nil
}

func (account *UserAccount) RemoveBookmark(suggestionId string) error {
	if _, ok := account.Bookmarks[suggestionId]; !ok {
		return ErrBookmarkNotFound
	}

	delete(account.Bookmarks, suggestionId)
	account.Update()

	return nil
}
