package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

//A Bookmark model
type Bookmark struct {
	Id           bson.ObjectId `json:"id"            bson:"_id"           binding:"required"`
	Creation     time.Time     `json:"creation"      bson:"creation"      binding:"required"`
	SuggestionId bson.ObjectId `json:"suggestion_id" bson:"suggestion_id" binding:"required"`
}

//Creates a new Bookmark for the specified suggestion. The creation will be set to time.now()
func NewBookmark(suggestionId bson.ObjectId) *Bookmark {
	return &Bookmark{
		Id:           bson.NewObjectId(),
		SuggestionId: suggestionId,
		Creation:     time.Now(),
	}
}
