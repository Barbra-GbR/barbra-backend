package models

import (
	"github.com/bitphinix/barbra-backend/payloads"
	"github.com/bitphinix/barbra-backend/db"
	"gopkg.in/mgo.v2/bson"
	"github.com/bitphinix/barbra-backend/helpers"
)

type UserProfile struct {
	Email               string `json:"email"       bson:"email,omitempty"       validate:"email,lowercase"              binding:"required"`
	GivenName           string `json:"given_name"  bson:"given_name,omitempty"  validate:"alphaunicode,min=1,max=50"    binding:"required"`
	FamilyName          string `json:"family_name" bson:"family_name,omitempty" validate:"alphaunicode,min=1,max=50"    binding:"required"`
	PictureURL          string `json:"picture"     bson:"picture,omitempty"     validate:"url"                          binding:"required"`

	//TODO: Unique                                                     .
	Nickname            string `json:"nickname"    bson:"nickname,omitempty"    validate:"alphanumunicode,min=2,max=50" binding:"required"`
}

func (profile *UserProfile) UpdateInfo(payload *payloads.ProfilePayload) error {
	collection := db.GetDB().C("users")
	validate := helpers.GetValidator()

	if err := validate.Struct(payload); err != nil {
		return ErrInvalidPayload
	}

	if profile.Email != payload.Email && payload.Email != "" {
		count, err := collection.Find(bson.M{"email": payload.Email}).Count()

		if err != nil || count > 0 {
			return ErrEmailAlreadyInUse
		}

		profile.Email = payload.Email
	}

	//The ugly but efficient way
	if payload.PictureURL != "" {
		profile.PictureURL = payload.PictureURL
	}

	if payload.Nickname != "" {
		profile.Nickname = payload.Nickname
	}

	if payload.GivenName != "" {
		profile.GivenName = payload.GivenName
	}

	if payload.FamilyName != "" {
		profile.FamilyName = payload.FamilyName
	}

	return nil
}
