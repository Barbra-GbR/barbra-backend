package auth

import (
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"context"
	"fmt"
	"errors"
	"encoding/base64"
	"crypto/rand"
	"github.com/bitphinix/babra_backend/config"
	"github.com/bitphinix/babra_backend/models"
)

var (
	ErrUnableToFetchOIDToken = errors.New("openIDProvider: Unable to fetch oidToken")
)

type OpenIDClient struct {
	oidVerifier  *oidc.IDTokenVerifier
	oidProvider  *oidc.Provider
	oAuth2Config *oauth2.Config
}

func LoadOpenIDClient(providerID string) (*OpenIDClient, error) {
	c := config.GetConfig()
	cAddr := "auth." + providerID

	oauth2Config := &oauth2.Config{
		ClientID:     c.GetString(cAddr + ".key"),
		ClientSecret: c.GetString(cAddr + ".secret"),
		RedirectURL:  getCallbackURL(c.GetString("server.host"), providerID),
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.GetString(cAddr + ".endpoint.auth_url"),
			TokenURL: c.GetString(cAddr + ".endpoint.token_url"),
		},
		Scopes: []string{
			oidc.ScopeOpenID,
			"email",
			"profile",
		},
	}

	oidProvider, err := oidc.NewProvider(context.Background(), c.GetString(cAddr + ".endpoint.url"))

	if err != nil {
		return nil, err
	}

	oidVerifier := oidProvider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID})

	return &OpenIDClient{
		oidProvider:  oidProvider,
		oAuth2Config: oauth2Config,
		oidVerifier:  oidVerifier,
	}, nil
}

func (client *OpenIDClient) GenerateLoginURL(state string) string {
	return client.oAuth2Config.AuthCodeURL(state)
}

func (client *OpenIDClient) FetchOAuthToken(code string) (*oauth2.Token, error) {
	return client.oAuth2Config.Exchange(context.Background(), code)
}

func (client *OpenIDClient) FetchOIDToken(oauth2Token *oauth2.Token) (*oidc.IDToken, error) {
	rawOIDToken, ok := oauth2Token.Extra("id_token").(string)

	if !ok {
		return nil, ErrUnableToFetchOIDToken
	}

	return client.oidVerifier.Verify(context.Background(), rawOIDToken)
}

func (client *OpenIDClient) GetUserSub(token *oidc.IDToken) string {
	return token.Subject
}

func (client *OpenIDClient) FetchUserInfo(token *oauth2.Token) (*models.UserInfo, error) {
	oidProfile, err := client.oidProvider.UserInfo(context.Background(), oauth2.StaticTokenSource(token))

	if err != nil {
		return nil, err
	}

	userInfo := new(models.UserInfo)
	err = oidProfile.Claims(userInfo)

	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func getCallbackURL(host string, providerID string) string {
	return fmt.Sprintf("%s/api/v1/login/%s/callback", host, providerID)
}

func GenerateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), err
}
