package auth

import (
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"context"
	"fmt"
	"errors"
	"encoding/base64"
	"crypto/rand"
	"github.com/Barbra-GbR/barbra-backend/config"
	"github.com/Barbra-GbR/barbra-backend/payloads"
)

var (
	ErrUnableToFetchOIdToken = errors.New("openIdProvider: Unable to fetch oidToken")
)

type OpenIdClient struct {
	oidVerifier  *oidc.IDTokenVerifier
	oidProvider  *oidc.Provider
	oAuth2Config *oauth2.Config
}

func LoadOpenIdClient(providerId string) (*OpenIdClient, error) {
	c := config.GetConfig()
	cAddr := "auth." + providerId

	oauth2Config := &oauth2.Config{
		ClientID:     c.GetString(cAddr + ".key"),
		ClientSecret: c.GetString(cAddr + ".secret"),
		RedirectURL:  getCallbackURL(c.GetString("server.host"), providerId),
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

	oidProvider, err := oidc.NewProvider(context.Background(), c.GetString(cAddr+".endpoint.url"))

	if err != nil {
		return nil, err
	}

	oidVerifier := oidProvider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID})

	return &OpenIdClient{
		oidProvider:  oidProvider,
		oAuth2Config: oauth2Config,
		oidVerifier:  oidVerifier,
	}, nil
}

func (client *OpenIdClient) GenerateLoginURL(state string) string {
	return client.oAuth2Config.AuthCodeURL(state)
}

func (client *OpenIdClient) FetchOAuthToken(code string) (*oauth2.Token, error) {
	return client.oAuth2Config.Exchange(context.Background(), code)
}

func (client *OpenIdClient) FetchOIdToken(oauth2Token *oauth2.Token) (*oidc.IDToken, error) {
	rawOIdToken, ok := oauth2Token.Extra("id_token").(string)

	if !ok {
		return nil, ErrUnableToFetchOIdToken
	}

	return client.oidVerifier.Verify(context.Background(), rawOIdToken)
}

func (client *OpenIdClient) GetUserSub(token *oidc.IDToken) string {
	return token.Subject
}

func (client *OpenIdClient) FetchProfilePayload(token *oauth2.Token) (*payloads.ProfilePayload, error) {
	oidProfile, err := client.oidProvider.UserInfo(context.Background(), oauth2.StaticTokenSource(token))

	if err != nil {
		return nil, err
	}

	payload := new(payloads.ProfilePayload)
	err = oidProfile.Claims(payload)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

func getCallbackURL(host string, providerId string) string {
	return fmt.Sprintf("%s/api/v1/login/%s/callback", host, providerId)
}

func GenerateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), err
}
