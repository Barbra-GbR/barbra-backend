package models

import (
	"github.com/Barbra-GbR/barbra-backend/db"
	"gopkg.in/mgo.v2/bson"
)

//Used for storing openId ids and assigning them to UserAccounts
type OpenIdAccount struct {
	Id       string        `json:"id" bson:"_id"`
	Provider string        `json:"provider" bson:"provider"`
	Sub      string        `json:"sub" bson:"sub"`
	OwnerId  bson.ObjectId `json:"owner" bson:"owner"`
}

//Returns the OpenIdAccount for the provider and subject
func FindOIdAccount(provider string, sub string) (*OpenIdAccount, error) {
	collection := db.GetDB().C("openid_accounts")
	account := new(OpenIdAccount)
	err := collection.FindId(getOIdAccountId(provider, sub)).One(account)
	return account, err
}

//Registers a new OpenIdAccount (This will override accounts!)
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

//Returns the accountId for the spevified provider and sub
func getOIdAccountId(provider string, sub string) string {
	return provider + "," + sub
}

//Saves the OpenIdAccount to the database
func (account *OpenIdAccount) Save() error {
	collection := db.GetDB().C("openid_accounts")
	_, err := collection.UpsertId(account.Id, account)
	return err
}
