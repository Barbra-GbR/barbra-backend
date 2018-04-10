package models

import (
	"github.com/Barbra-GbR/barbra-backend/db"
	"gopkg.in/mgo.v2/bson"
)

type OpenIdAccount struct {
	Id       string        `json:"id" bson:"_id"`
	Provider string        `json:"provider" bson:"provider"`
	Sub      string        `json:"sub" bson:"sub"`
	OwnerId  bson.ObjectId `json:"owner" bson:"owner"`
}

func FindOIdAccount(provider string, sub string) (*OpenIdAccount, error) {
	collection := db.GetDB().C("openid_accounts")

	account := new(OpenIdAccount)
	err := collection.FindId(getOIdAccountId(provider, sub)).One(account)
	return account, err
}

func RegisterOIdAccount(provider string, sub string, ownerId bson.ObjectId) (*OpenIdAccount, error) {
	account := &OpenIdAccount{
		Provider: provider,
		Sub:      sub,
		Id:       getOIdAccountId(provider, sub),
		OwnerId:  ownerId,
	}

	err := account.Save()
	return account, err
}

func getOIdAccountId(provider string, sub string) string {
	return provider + "," + sub
}

func (account *OpenIdAccount) Save() error {
	collection := db.GetDB().C("openid_accounts")
	_, err := collection.UpsertId(account.Id, account)
	return err
}
