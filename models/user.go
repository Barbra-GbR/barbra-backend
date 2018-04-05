package models

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"../db"
)

var (
	ErrInvalidUserInfo   = errors.New("user: invalid UserInfo")
	ErrEmailAlreadyInUse = errors.New("user: email already in use")
)

type UserAccount struct {
	*UserInfo `bson:"user_info"`
	ID string `json:"id" bson:"_id"`
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

func RegisterUser(userInfo *UserInfo) (*UserAccount, error) {
	collection := db.GetDB().C("users")

	userInfo.Normalize()

	if !userInfo.Verify() {
		return nil, ErrInvalidUserInfo
	}

	count, err := collection.Find(bson.M{"user_info.email": userInfo.Email}).Count()

	if err != nil || count > 0 {
		return nil, ErrEmailAlreadyInUse
	}

	account := &UserAccount{
		UserInfo: userInfo,
		ID:       bson.NewObjectId().Hex(),
	}

	err = collection.Insert(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (account *UserAccount) UpdateAccountInfo(userInfo *UserInfo) error {
	collection := db.GetDB().C("users")

	if !userInfo.Verify() {
		return ErrInvalidUserInfo
	}

	count, err := collection.Find(bson.M{"user_info.email": userInfo.Email}).Count()

	if err != nil || count > 0 {
		return ErrEmailAlreadyInUse
	}

	account.UserInfo = userInfo
	return account.Save()
}

func (account *UserAccount) Save() error {
	collection := db.GetDB().C("users")
	return collection.UpdateId(account.ID, account)
}
