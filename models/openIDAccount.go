package models

import (
	"github.com/bitphinix/barbra_backend/db"
)

type OpenIdAccount struct {
	Id       string `json:"id" bson:"_id"`
	Provider string `json:"provider" bson:"provider"`
	Sub      string `json:"sub" bson:"sub"`
	Owner    string `json:"owner" bson:"owner"`
}

func FindOIdAccount(provider string, sub string) (*OpenIdAccount, error) {
	collection := db.GetDB().C("openid_accounts")

	account := new(OpenIdAccount)
	err := collection.FindId(getOIdAccountId(provider, sub)).One(account)
	return account, err
}

func RegisterOIdAccount(provider string, sub string, owner string) (*OpenIdAccount, error) {
	account := &OpenIdAccount{
		Provider: provider,
		Sub:      sub,
		Id:       getOIdAccountId(provider, sub),
		Owner:    owner,
	}

	err := account.Save()
	return account, err
}

func getOIdAccountId(provider string, sub string) string {
	return provider + "," + sub
}

func (account *OpenIdAccount) Save() error {
	collection := db.GetDB().C("openid_accounts")
	return collection.Insert(account)
}

func (account *OpenIdAccount) Update() error {
	collection := db.GetDB().C("openid_accounts")
	return collection.UpdateId(account.Id, account)
}
