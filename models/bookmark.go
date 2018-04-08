package models

import (
	"time"
	"errors"
	"github.com/bitphinix/barbra-backend/db"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrBookmarkExists   = errors.New("bookmark: already exists")
)

type Bookmark struct {
	Id           string    `json:"id"            bson:"_id"           validate:"hexadecimal" binding:"required"`
	Creation     time.Time `json:"creation"      bson:"creation"                             binding:"required"`
	SuggestionId string    `json:"suggestion_id" bson:"suggestion_id" validate:"hexadecimal" binding:"required"`
}

func NewBookmark(suggestionId string) *Bookmark {
	return &Bookmark{
		SuggestionId: suggestionId,
		Creation:     time.Now(),
	}
}

type BookmarkContainer struct {
	Id string `json:"id"        bson:"_id"       validate:"hexadecimal" binding:"required"`
}

func NewBookmarkContainer() (*BookmarkContainer, error) {
	collection := db.GetDB().C("bookmark_container")

	container := &BookmarkContainer{
		Id: bson.NewObjectId().Hex(),
	}

	return container, collection.Insert(bson.M{"_id": container.Id, "bookmarks": nil})
}

func FindBookmarkContainerId(id string) (*BookmarkContainer, error) {
	collection := db.GetDB().C("bookmark_container")
	container := new(BookmarkContainer)
	err := collection.FindId(id).One(container)
	return container, err
}

func (container *BookmarkContainer) AddBookmark(suggestionId string) error {
	collection := db.GetDB().C("bookmark_container")

	if !SuggestionExists(suggestionId) {
		return ErrSuggestionNotFound
	}

	if container.ContainsBookmark(suggestionId) {
		return ErrBookmarkExists
	}

	return collection.Update(container.Id, bson.M{"bookmarks": bson.M{"$push": NewBookmark(suggestionId)}})
}

func (container *BookmarkContainer) ContainsBookmark(suggestionId string) bool {
	collection := db.GetDB().C("bookmark_container")
	count, err := collection.Find(bson.M{"_id": container.Id, "bookmarks": bson.M{"$in": bson.M{"suggestionId": suggestionId}}}).Limit(1).Count()
	return count == 1 && err == nil
}

func (container *BookmarkContainer) RemoveBookmark(suggestionId string) error {
	collection := db.GetDB().C("bookmark_container")
	return collection.Update(container.Id, bson.M{"bookmarks": bson.M{"$pull": bson.M{"suggestionId": suggestionId}}})
}

