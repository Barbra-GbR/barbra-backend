package models

import (
	"time"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrBookmarkExists = errors.New("bookmark: already exists")
)

type Bookmark struct {
	Id           bson.ObjectId `json:"id"            bson:"_id"           binding:"required"`
	Creation     time.Time     `json:"creation"      bson:"creation"      binding:"required"`
	SuggestionId bson.ObjectId `json:"suggestion_id" bson:"suggestion_id" binding:"required"`
}

func NewBookmark(suggestionId bson.ObjectId) *Bookmark {
	return &Bookmark{
		Id:           bson.NewObjectId(),
		SuggestionId: suggestionId,
		Creation:     time.Now(),
	}
}
