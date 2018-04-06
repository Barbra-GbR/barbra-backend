package models

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/bitphinix/barbra_backend/db"
)

var (
	ErrInvalidUserInfo   = errors.New("user: invalid UserInfo")
	ErrEmailAlreadyInUse = errors.New("user: email already in use")
)

type UserAccount struct {
	*UserInfo `bson:"user_info"`
	Id string `json:"id" bson:"_id"`
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
		Id:       bson.NewObjectId().Hex(),
	}

	err = collection.Insert(account)

	return account, err
}

func (account *UserAccount) UpdateAccountInfo(userInfo *UserInfo) error {
	collection := db.GetDB().C("users")

	userInfo.Normalize()

	if !userInfo.Verify() {
		return ErrInvalidUserInfo
	}

	if account.Email != userInfo.Email {
		count, err := collection.Find(bson.M{"user_info.email": userInfo.Email}).Count()

		if err != nil || count > 0 {
			return ErrEmailAlreadyInUse
		}
	}

	account.UserInfo = userInfo
	return account.Update()
}

func (account *UserAccount) Update() error {
	collection := db.GetDB().C("users")
	return collection.UpdateId(account.Id, account)
}

func (account *UserAccount) Delete() error {
	collection := db.GetDB().C("users")
	return collection.RemoveId(account.Id)
}
