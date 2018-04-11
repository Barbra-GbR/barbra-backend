package models

import (
	"github.com/Barbra-GbR/barbra-backend/db"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"errors"
)

var (
	ErrSuggestionNotFound = errors.New("suggestion: not found")
)

//A Suggestion model
type Suggestion struct {
	Provider string        `json:"provider" bson:"provider" binding:"required"`
	Url      string        `json:"url"      bson:"url"      binding:"required"`
	Kind     string        `json:"kind"     bson:"kind"     binding:"required"`
	Title    string        `json:"title"    bson:"title"    binding:"required"`
	Category string        `json:"category" bson:"category" binding:"required"`
	Tags     []string      `json:"tags"     bson:"tags"     binding:"required"`
	Content  string        `json:"content"  bson:"content"  binding:"required"`
	Id       bson.ObjectId `json:"id"       bson:"_id"      binding:"required"`
}

//Creates a new suggestion with the specified data
func NewSuggestion(url string, kind string, title string, category string, provider string, tags []string, content string) *Suggestion {
	return &Suggestion{
		Id:       bson.NewObjectId(),
		Content:  content,
		Kind:     kind,
		Tags:     tags,
		Title:    title,
		Category: category,
		Url:      url,
		Provider: provider,
	}
}

//Searches for the suggestion in the database, if no matches are found a new one will be created
func GetSuggestion(url string, kind string, title string, provider string, category string, tags []string, content string) (*Suggestion, error) {
	collection := db.GetDB().C("suggestions")

	suggestion := new(Suggestion)
	err := collection.Find(bson.M{
		"url":      url,
		"kind":     kind,
		"title":    title,
		"category": category,
		"tags":     tags,
		"content":  content,
		"provider": provider,
	}).One(suggestion)
	if err == mgo.ErrNotFound {
		suggestion = NewSuggestion(url, kind, title, category, provider, tags, content)
		err = suggestion.Save()
	}

	return suggestion, err
}

//Returns the suggestion with the matching id
func GetSuggestionById(id bson.ObjectId) (*Suggestion, error) {
	collection := db.GetDB().C("suggestions")

	suggestion := new(Suggestion)
	err := collection.FindId(id).One(suggestion)
	return suggestion, err
}

//Checks if the specified suggestion exists
func SuggestionExists(id bson.ObjectId) bool {
	collection := db.GetDB().C("suggestions")
	count, err := collection.FindId(id).Limit(1).Count()
	return count > 0 && err == nil
}

//Saves the suggestion to the database
func (suggestion *Suggestion) Save() error {
	collection := db.GetDB().C("suggestions")
	_, err := collection.UpsertId(suggestion.Id, suggestion)
	return err
}
