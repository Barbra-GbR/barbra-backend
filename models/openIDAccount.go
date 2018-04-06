package models

import (
	"gopkg.in/mgo.v2"
	"github.com/bitphinix/babra_backend/db"
)

type OpenIDAccount struct {
	Id       string `json:"id" bson:"_id"`
	Provider string `json:"provider" bson:"provider"`
	Sub      string `json:"sub" bson:"sub"`
	Owner    string `json:"owner" bson:"owner"`
}

func GetOIDAccount(provider string, sub string) (*OpenIDAccount, error) {
	collection := db.GetDB().C("openid_accounts")

	account := new(OpenIDAccount)
	err := collection.FindId(getOIDAccountID(provider, sub)).One(account)

	if err != nil {
		return nil, err
	}

	return account, err
}

func RegisterOIDAccount(provider string, sub string, owner string) (*OpenIDAccount, error) {
	account := &OpenIDAccount{
		Provider: provider,
		Sub:      sub,
		Id:       getOIDAccountID(provider, sub),
		Owner:    owner,
	}

	err := account.Save()

	if err != nil {
		return nil, err
	}

	return account, nil
}

func getOIDAccountID(provider string, sub string) string {
	return provider + "," + sub
}

func (account *OpenIDAccount) Save() error {
	collection := db.GetDB().C("openid_accounts")

	err := collection.UpdateId(account.Id, account)

	if err == mgo.ErrNotFound {
		err = collection.Insert(account)
	}

	return err
}
