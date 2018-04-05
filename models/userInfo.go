package models

import (
	"strings"
	"regexp"
)

type UserInfo struct {
	Email      string `json:"email" bson:"email" binding:"required"`
	GivenName  string `json:"given_name" bson:"given_name" binding:"required"`
	FamilyName string `json:"family_name" bson:"family_name" binding:"required"`
}

func (userInfo *UserInfo) Normalize() {
	userInfo.Email = strings.ToLower(userInfo.Email)
}

func (userInfo *UserInfo) Verify() bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(userInfo.Email) {
		return false
	}

	if len(userInfo.GivenName) < 1 {
		return false
	}

	return len(userInfo.FamilyName) > 0
}
