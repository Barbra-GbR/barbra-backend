package models

import (
	"github.com/Barbra-GbR/barbra-backend/db"
	"github.com/Barbra-GbR/barbra-backend/helpers"
	"github.com/Barbra-GbR/barbra-backend/payloads"
	"github.com/juju/errors"
)

//The UserProfile model
type UserProfile struct {
	Email      string `json:"email"       	bson:"email,omitempty"     		validate:"email,lowercase"           		binding:"required"`
	GivenName  string `json:"given_name"  	bson:"given_name,omitempty"  	validate:"alphaunicode,min=1,max=50"		binding:"required"`
	FamilyName string `json:"family_name" 	bson:"family_name,omitempty" 	validate:"alphaunicode,min=1,max=50" 		binding:"required"`
	PictureURL string `json:"picture"     	bson:"picture,omitempty"     	validate:"url"                       		binding:"required"`
	Nickname   string `json:"nickname"   	bson:"nickname,omitempty"    	validate:"alphanumunicode,min=2,max=50" 	binding:"required"` //TODO: Unique
}

// UpdateInfo updates the profile with the given payload and validates it. Empty fields won't be set
func (profile *UserProfile) UpdateInfo(payload *payloads.ProfilePayload) error {
	collection := db.GetDB().C("users")
	validate := helpers.GetValidator()

	if err := validate.Struct(payload); err != nil {
		return errors.NewNotValid(err, "invalid payload")
	}

	if profile.Email != payload.Email && payload.Email != "" {
		if UserEmailInUse(payload.Email) {
			return errors.AlreadyExistsf("the email %s already is already taken", payload.Email)
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
