package auth

import (
	"errors"
	"gopkg.in/mgo.v2"
	"golang.org/x/oauth2"
	"github.com/coreos/go-oidc"
	"github.com/bitphinix/barbra_backend/config"
	"github.com/bitphinix/barbra_backend/models"
)

var (
	ErrProviderAlreadyRegistered = errors.New("openIdProvider: Provider is already registered")
	ErrClientNotFound            = errors.New("openIdProvider: No client with the specified providerId registered")
	userManager                  *UserManager
)

type UserManager struct {
	oidClients map[string]*OpenIdClient
}

func GetUserManager() *UserManager {
	return userManager
}

func InitUserManager() {
	c := config.GetConfig()
	manager := new(UserManager)
	manager.oidClients = make(map[string]*OpenIdClient)

	for providerId := range c.GetStringMap("auth") {
		err := manager.LoadOIdProvider(providerId)

		if err != nil {
			panic(err)
		}
	}

	userManager = manager
}

func (manager *UserManager) LoadOIdProvider(providerId string) error {
	if _, ok := manager.oidClients[providerId]; ok {
		return ErrProviderAlreadyRegistered
	}

	client, err := LoadOpenIdClient(providerId)

	if err != nil {
		return err
	}

	manager.oidClients[providerId] = client
	return nil
}

func (manager *UserManager) GetOIdClient(providerId string) (*OpenIdClient, error) {
	var client *OpenIdClient
	var ok bool

	if client, ok = manager.oidClients[providerId]; !ok {
		return nil, ErrClientNotFound
	}

	return client, nil
}

func (manager *UserManager) GenerateLoginUrl(providerId string, state string) (string, error) {
	client, err := manager.GetOIdClient(providerId)
	if err != nil {
		return "", err
	}

	return client.GenerateLoginURL(state), nil
}

func (manager *UserManager) GetAccount(providerId string, code string) (*models.UserAccount, error) {
	client, err := manager.GetOIdClient(providerId)
	if err != nil {
		return nil, err
	}

	oauthToken, err := client.FetchOAuthToken(code)
	if err != nil {
		return nil, err
	}

	oidToken, err := client.FetchOIdToken(oauthToken)
	if err != nil {
		return nil, err
	}

	account, err := manager.FindAccount(providerId, oidToken)
	if err != mgo.ErrNotFound {
		return account, err
	}

	return manager.RegisterAccount(providerId, oauthToken, oidToken)
}

func (manager *UserManager) RegisterAccount(providerId string, oauth2Token *oauth2.Token, oidToken *oidc.IDToken) (*models.UserAccount, error) {
	client, err := manager.GetOIdClient(providerId)
	if err != nil {
		return nil, err
	}

	userInfo, err := client.FetchUserInfo(oauth2Token)
	if err != nil {
		return nil, err
	}

	account, err := models.RegisterUser(userInfo)
	if err != nil {
		return nil, err
	}

	_, err = models.RegisterOIdAccount(providerId, oidToken.Subject, account.Id)

	if err != nil {
		_ = account.Delete()
		
		return nil, err
	}

	return account, nil
}

func (UserManager) FindAccount(providerId string, oidToken *oidc.IDToken) (*models.UserAccount, error) {
	oidAccount, err := models.FindOIdAccount(providerId, oidToken.Subject)
	if err != nil {
		return nil, err
	}

	account, err := models.GetUserAccount(oidAccount.Owner)
	return account, err
}
