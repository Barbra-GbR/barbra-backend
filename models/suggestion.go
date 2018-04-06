package models

import (
	"github.com/bitphinix/barbra_backend/db"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type Suggestion struct {
	Url     string   `json:"url" bson:"url" binding:"required"`
	Kind    string   `json:"kind" bson:"kind" binding:"required"`
	Title   string   `json:"title" bson:"title" binding:"required"`
	Topic   string   `json:"topic" bson:"topic" binding:"required"`
	Tags    []string `json:"tags" bson:"tags" binding:"required"`
	Content string   `json:"content" bson:"content" binding:"required"`
	LogoUrl string   `json:"logo_url" bson:"logo_url" binding:"required"`
	Id      string   `json:"id" bson:"_id" binding:"required"`
}

func NewSuggestion(url string, kind string, title string, topic string, tags []string, content string, logoUrl string) *Suggestion {
	return &Suggestion{
		Id:      bson.NewObjectId().Hex(),
		Content: content,
		Kind:    kind,
		LogoUrl: logoUrl,
		Tags:    tags,
		Title:   title,
		Topic:   topic,
		Url:     url,
	}
}

func GetSuggestion(url string, kind string, title string, topic string, tags []string, content string, logoUrl string) (*Suggestion, error) {
	collection := db.GetDB().C("suggestions")

	suggestion := new(Suggestion)
	err := collection.Find(bson.M{
		"url":      url,
		"kind":     kind,
		"title":    title,
		"topic":    topic,
		"tags":     tags,
		"content":  content,
		"logo_url": logoUrl,
	}).One(suggestion)

	if err == mgo.ErrNotFound {
		suggestion := NewSuggestion(url, kind, title, topic, tags, content, logoUrl)
		err = suggestion.Save()
	}

	return suggestion, err
}

func FindSuggestionById(id string) (*Suggestion, error) {
	collection := db.GetDB().C("suggestions")

	suggestion := new(Suggestion)
	err := collection.FindId(id).One(suggestion)
	return suggestion, err
}

func (suggestion Suggestion) Save() error {
	collection := db.GetDB().C("suggestions")
	return collection.Insert(suggestion)
}

func (suggestion *Suggestion) Update() error {
	collection := db.GetDB().C("suggestions")
	return collection.UpdateId(suggestion.Id, suggestion)
}
