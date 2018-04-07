package models

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/bitphinix/barbra_backend/db"
	"github.com/bitphinix/barbra_backend/helpers"
	"github.com/bitphinix/barbra_backend/payloads"
	"log"
)

var (
	ErrInvalidUserInfo   = errors.New("user: invalid UserInfo")
	ErrEmailAlreadyInUse = errors.New("user: email already in use")
)

type UserAccount struct {
	Id         string `json:"id"          bson:"_id"                   validate:"hexadecimal"                  binding:"required"`
	Enrolled   bool   `json:"enrolled"    bson:"enrolled"                                                      binding:"required"`
	Email      string `json:"email"       bson:"email,omitempty"       validate:"email,lowercase"              binding:"required"`
	GivenName  string `json:"given_name"  bson:"given_name,omitempty"  validate:"alphaunicode,min=1,max=50"    binding:"required"`
	FamilyName string `json:"family_name" bson:"family_name,omitempty" validate:"alphaunicode,min=1,max=50"    binding:"required"`
	PictureURL string `json:"picture"     bson:"picture,omitempty"     validate:"url"                          binding:"required"`
	Nickname   string `json:"nickname"    bson:"nickname,omitempty"    validate:"alphanumunicode,min=2,max=50" binding:"required"`
}

func GetUserAccount(id string) (*UserAccount, error) {
	collection := db.GetDB().C("users")

	account := new(UserAccount)
	err := collection.FindId(id).One(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (account *UserAccount) IsEnrolled() bool {
	validate := helpers.GetValidator()
	err := validate.Struct(account)
	return err == nil
}

func RegisterUser(payload *payloads.ProfilePayload) (*UserAccount, error) {
	collection := db.GetDB().C("users")
	validate := helpers.GetValidator()

	if err := validate.Struct(payload); err != nil {
		log.Println(err.Error())
		return nil, ErrInvalidUserInfo
	}

	//Check if email is already in use
	if payload.Email != "" {
		count, err := collection.Find(bson.M{"email": payload.Email}).Count()
		if err != nil || count > 0 {
			return nil, ErrEmailAlreadyInUse
		}
	}

	account := &UserAccount{
		Email:      payload.Email,
		Id:         bson.NewObjectId().Hex(),
		Enrolled:   false,
		FamilyName: payload.FamilyName,
		GivenName:  payload.GivenName,
		Nickname:   payload.Nickname,
		PictureURL: payload.PictureURL,
	}

	account.Enrolled = account.IsEnrolled()
	return account, account.Save()
}

func (account *UserAccount) UpdateAccountInfo(payload *payloads.ProfilePayload) error {
	collection := db.GetDB().C("users")
	validate := helpers.GetValidator()

	if err := validate.Struct(payload); err != nil {
		return ErrInvalidUserInfo
	}

	if account.Email != payload.Email && payload.Email != "" {
		count, err := collection.Find(bson.M{"email": payload.Email}).Count()

		if err != nil || count > 0 {
			return ErrEmailAlreadyInUse
		}

		account.Email = payload.Email
	}

	//The ugly but efficient way
	if payload.PictureURL != "" {
		account.PictureURL = payload.PictureURL
	}
	if payload.Nickname != "" {
		account.Nickname = payload.Nickname
	}
	if payload.GivenName != "" {
		account.GivenName = payload.GivenName
	}
	if payload.FamilyName != "" {
		account.FamilyName = payload.FamilyName
	}

	account.Enrolled = account.IsEnrolled()
	return account.Update()
}

func (account *UserAccount) Save() error {
	collection := db.GetDB().C("users")
	return collection.Insert(account)
}

func (account *UserAccount) Update() error {
	collection := db.GetDB().C("users")
	return collection.UpdateId(account.Id, account)
}

func (account *UserAccount) Delete() error {
	collection := db.GetDB().C("users")
	return collection.RemoveId(account.Id)
}
