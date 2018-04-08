package models

import (
	"time"
	"errors"
	"github.com/bitphinix/barbra_backend/db"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrBookmarkExists   = errors.New("bookmark: already exists")
	ErrBookmarkNotFound = errors.New("bookmark: bookmark not found")
)

type Bookmark struct {
	Id           string    `json:"id"            bson:"_id"           validate:"hexadecimal" binding:"required"`
	Creation     time.Time `json:"creation"      bson:"creation"                             binding:"required"`
	SuggestionId string    `json:"suggestion_id" bson:"suggestion_id" validate:"hexadecimal" binding:"required"`
}

type BookmarkContainer struct {
	Id           string   `json:"id"        bson:"_id"       validate:"hexadecimal" binding:"required"`
	Bookmarks []*Bookmark `json:"bookmarks" bson:"bookmarks" binding:"required"`
}

func FindBookmarkContainerId(id string) (*BookmarkContainer, error) {
	collection := db.GetDB().C("bookmarks")
	container := new(BookmarkContainer)
	err := collection.FindId(id).One(container)
	return container, err
}

func NewContainer() *BookmarkContainer {
	return &BookmarkContainer{
		Id:bson.NewObjectId().Hex(),
		Bookmarks:nil,
	}
}

func (container *BookmarkContainer) AddBookmark(suggestionId string) error {
	if _, ok := account.Bookmarks[suggestionId]; ok {
		return ErrBookmarkExists
	}

	if !SuggestionExists(suggestionId) {
		return ErrSuggestionNotFound
	}

	bookmark := NewBookmark(suggestionId)
	account.Bookmarks[suggestionId] = bookmark
	account.Update()

	return nil
}

func (container *BookmarkContainer) RemoveBookmark(suggestionId string) error {
	if _, ok := account.Bookmarks[suggestionId]; !ok {
		return ErrBookmarkNotFound
	}

	delete(account.Bookmarks, suggestionId)
	account.Update()

	return nil
}

func (container *BookmarkContainer) Save() error {
	collection := db.GetDB().C("bookmarks")
	return collection.Insert(container)
}

func (container *BookmarkContainer) Update() error {
	collection := db.GetDB().C("bookmarks")
	return collection.UpdateId(container.Id, container)
}

func (container *BookmarkContainer) Delete() error {
	collection := db.GetDB().C("bookmarks")
	return collection.RemoveId(container.Id)
}


func NewBookmark(suggestionId string) *Bookmark {
	return &Bookmark{
		SuggestionId: suggestionId,
		Creation:     time.Now(),
	}
}
