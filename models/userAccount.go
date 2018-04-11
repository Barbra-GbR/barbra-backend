package models

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/Barbra-GbR/barbra-backend/db"
	"github.com/Barbra-GbR/barbra-backend/helpers"
	"github.com/Barbra-GbR/barbra-backend/payloads"
	"log"
)

var (
	ErrInvalidPayload    = errors.New("user: invalid payload")
	ErrEmailAlreadyInUse = errors.New("user: email already in use")
)

//The UserAccount model
type UserAccount struct {
	Id                  bson.ObjectId `json:"id"       bson:"_id"                   binding:"required"`
	Enrolled            bool          `json:"enrolled" bson:"enrolled"              binding:"required"`
	Profile             *UserProfile  `json:"profile"  bson:"profile"               binding:"required"`
	BookmarkContainerId bson.ObjectId `json:"-"        bson:"bookmark_container_id" binding:"required"`
}

//Returns the UserAccount with the specified id
func GetUserAccountById(id bson.ObjectId) (*UserAccount, error) {
	collection := db.GetDB().C("users")

	account := new(UserAccount)
	err := collection.FindId(id).One(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

//Registers a new user with the specified payload and validates it
func RegisterUser(payload *payloads.ProfilePayload) (*UserAccount, error) {
	collection := db.GetDB().C("users")
	validate := helpers.GetValidator()

	//Check if payload is valid
	if err := validate.Struct(payload); err != nil {
		log.Println(err.Error())
		return nil, ErrInvalidPayload
	}

	//Check if email is already in use
	if payload.Email != "" {
		count, err := collection.Find(bson.M{"email": payload.Email}).Count()
		if err != nil || count > 0 {
			return nil, ErrEmailAlreadyInUse
		}
	}

	bookmarkContainer, err := NewBookmarkContainer()
	if err != nil {
		return nil, err
	}

	account := &UserAccount{
		Id:       bson.NewObjectId(),
		Enrolled: false,
		Profile: &UserProfile{
			Email:      payload.Email,
			FamilyName: payload.FamilyName,
			GivenName:  payload.GivenName,
			Nickname:   payload.Nickname,
			PictureURL: payload.PictureURL,
		},
		BookmarkContainerId: bookmarkContainer.Id,
	}

	account.Enrolled = account.IsEnrolled()
	return account, account.Save()
}

//Checks if the user has completed profile
func (account *UserAccount) IsEnrolled() bool {
	validate := helpers.GetValidator()
	err := validate.Struct(account)
	return err == nil
}

//Updates the UserProfile with the specifed info and validates it
func (account *UserAccount) UpdateProfile(payload *payloads.ProfilePayload) error {
	err := account.Profile.UpdateInfo(payload)

	if err != nil {
		return err
	}

	account.Enrolled = account.IsEnrolled()
	return account.Save()
}

//Returns the users bookmark-container
func (account *UserAccount) GetBookmarkContainer() (*BookmarkContainer, error) {
	return GetBookmarkContainerById(account.BookmarkContainerId)
}

//Saves the UserAccount to the database
func (account *UserAccount) Save() error {
	collection := db.GetDB().C("users")
	_, err := collection.UpsertId(account.Id, account)
	return err
}

//Deletes the useraccount from the database TODO Remove bookmarks etc too
func (account *UserAccount) Delete() error {
	collection := db.GetDB().C("users")
	return collection.RemoveId(account.Id)
}
