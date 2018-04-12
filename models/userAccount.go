package models

import (
	"github.com/Barbra-GbR/barbra-backend/db"
	"github.com/Barbra-GbR/barbra-backend/helpers"
	"github.com/Barbra-GbR/barbra-backend/payloads"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

var (
	userAccountCollectionName = "users"
)

// The UserAccount model
type UserAccount struct {
	ID                  bson.ObjectId `json:"id"       bson:"_id"                   binding:"required"`
	Enrolled            bool          `json:"enrolled" bson:"enrolled"              binding:"required"`
	Profile             *UserProfile  `json:"profile"  bson:"profile"               binding:"required"`
	BookmarkContainerID bson.ObjectId `json:"-"        bson:"bookmark_container_id" binding:"required"`
}

// GetUserAccountByID returns a UserAccount with the specified id from the Database
func GetUserAccountByID(id bson.ObjectId) (*UserAccount, error) {
	collection := db.GetDB().C("users")
	account := new(UserAccount)
	err := collection.FindId(id).One(account)
	return account, err
}

// RegisterUser crates a new user with the specified payload and saves him to the Database
// It checks for duplicate emails in the Database and validates the payload
func RegisterUser(payload *payloads.ProfilePayload) (*UserAccount, error) {
	collection := db.GetDB().C(userAccountCollectionName)
	validate := helpers.GetValidator()

	//Check if payload is valid
	if err := validate.Struct(payload); err != nil {
		return nil, errors.NewNotValid(err, "invalid profile payload")
	}

	//Check if email is already in use
	if payload.Email != "" && UserEmailInUse(payload.Email) {
		return nil, errors.AlreadyExistsf("the email %s already is already taken", payload.Email)
	}

	bookmarkContainer, err := NewBookmarkContainer()
	if err != nil {
		return nil, err
	}

	account := &UserAccount{
		ID:       bson.NewObjectId(),
		Enrolled: false,
		Profile: &UserProfile{
			Email:      payload.Email,
			FamilyName: payload.FamilyName,
			GivenName:  payload.GivenName,
			Nickname:   payload.Nickname,
			PictureURL: payload.PictureURL,
		},
		BookmarkContainerID: bookmarkContainer.ID,
	}

	account.Enrolled = account.IsEnrolled()
	return account, account.Save()
}

// UserEmailInUse checks if the given email is already exists in the Database
// The email string should be lowercase and trimmed. It won´t be validated
// Returns false if any errors occur
func UserEmailInUse(email string) bool {
	collection := db.GetDB().C(userAccountCollectionName)
	count, err := collection.Find(bson.M{"email": email}).Count()
	return count > 0 && err == nil
}

// IsEnrolled checks if the user has completed profile.
// If errors occur false will be returned
func (account *UserAccount) IsEnrolled() bool {
	validate := helpers.GetValidator()
	err := validate.Struct(account)
	return err == nil
}

// UpdateProfile updates the UserProfile and saves it to the Database with the specified info and validates it
// It checks for duplicate emails in the Database and validates the payload
func (account *UserAccount) UpdateProfile(payload *payloads.ProfilePayload) error {
	err := account.Profile.UpdateInfo(payload)
	if err != nil {
		return err
	}

	account.Enrolled = account.IsEnrolled()
	return account.Save()
}

// GetBookmarkContainer returns the users bookmark-container
// TODO: Without database call
func (account *UserAccount) GetBookmarkContainer() (*BookmarkContainer, error) {
	return GetBookmarkContainerById(account.BookmarkContainerID)
}

// Save inserts the account into the database. If it already exists it will be updated
func (account *UserAccount) Save() error {
	collection := db.GetDB().C(userAccountCollectionName)
	_, err := collection.UpsertId(account.ID, account)
	return err
}

// Delete deletes the useraccount from the database TODO: Remove bookmarks etc too
func (account *UserAccount) Delete() error {
	collection := db.GetDB().C(userAccountCollectionName)
	return collection.RemoveId(account.ID)
}
