package models

import (
	"strings"
	"html"
	"github.com/bitphinix/barbra_backend/helpers"
)

type UserInfo struct {
	Email      string `json:"email" bson:"email" binding:"required"`
	GivenName  string `json:"given_name" bson:"given_name" binding:"required"`
	FamilyName string `json:"family_name" bson:"family_name" binding:"required"`
	PictureURL string `json:"picture" bson:"picture" binding:"required"`
	Nickname   string `json:"nickname" bson:"nickname" binding:"required"`
}

func (userInfo *UserInfo) Normalize() {
	userInfo.Email = html.EscapeString(strings.ToLower(userInfo.Email))
	userInfo.FamilyName = html.EscapeString(userInfo.FamilyName)
	userInfo.GivenName = html.EscapeString(userInfo.GivenName)
}

func (userInfo *UserInfo) Verify() bool {
	if !helpers.RegexEmail.MatchString(userInfo.Email) {
		return false
	}

	if !helpers.RegexName.MatchString(userInfo.GivenName) || len(userInfo.GivenName) > 50 || len(userInfo.GivenName) < 1 {
		return false
	}

	if !helpers.RegexNickname.MatchString(userInfo.Nickname) || len(userInfo.Nickname) > 50 || len(userInfo.Nickname) < 3 {
		return false
	}

	if !helpers.RegexURL.MatchString(userInfo.PictureURL) || len(userInfo.PictureURL) > 100 || len(userInfo.PictureURL) < 1 {
		return false
	}

	return helpers.RegexName.MatchString(userInfo.FamilyName) && len(userInfo.FamilyName) >= 1 && len(userInfo.FamilyName) <= 50
}
