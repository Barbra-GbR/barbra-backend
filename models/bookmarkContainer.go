package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/Barbra-GbR/barbra-backend/db"
)

type BookmarkContainer struct {
	Id bson.ObjectId `json:"id" bson:"_id" binding:"required"`
}

func NewBookmarkContainer() (*BookmarkContainer, error) {
	collection := db.GetDB().C("bookmark_container")

	container := &BookmarkContainer{
		Id: bson.NewObjectId(),
	}

	return container, collection.Insert(container)
}

func GetBookmarkContainerById(id bson.ObjectId) (*BookmarkContainer, error) {
	collection := db.GetDB().C("bookmark_container")
	container := new(BookmarkContainer)
	err := collection.FindId(id).One(container)
	return container, err
}

func (container *BookmarkContainer) AddBookmark(suggestionId bson.ObjectId) error {
	collection := db.GetDB().C("bookmark_container")

	if !SuggestionExists(suggestionId) {
		return ErrSuggestionNotFound
	}

	if container.ContainsBookmark(suggestionId) {
		return ErrBookmarkExists
	}

	_, err := collection.UpsertId(container.Id, bson.M{"$push": bson.M{"bookmarks": NewBookmark(suggestionId)}})
	return err
}

func (container *BookmarkContainer) ContainsBookmark(suggestionId bson.ObjectId) bool {
	collection := db.GetDB().C("bookmark_container")
	count, err := collection.Find(bson.M{"_id": container.Id, "$in": bson.M{"bookmarks": bson.M{"suggestionId": suggestionId}}}).Limit(1).Count()
	return count == 1 && err == nil
}

func (container *BookmarkContainer) RemoveBookmark(suggestionId bson.ObjectId) error {
	collection := db.GetDB().C("bookmark_container")
	return collection.UpdateId(container.Id, bson.M{"$pull": bson.M{"bookmarks": bson.M{"suggestionId": suggestionId}}})
}
