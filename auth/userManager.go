package auth

import (
	"errors"
	"gopkg.in/mgo.v2"
	"golang.org/x/oauth2"
	"github.com/coreos/go-oidc"
	"github.com/bitphinix/babra_backend/config"
	"github.com/bitphinix/babra_backend/models"
)

var (
	ErrProviderAlreadyRegistered = errors.New("openIDProvider: Provider is already registered")
	ErrClientNotFound            = errors.New("openIDProvider: No client with the specified providerID registered")
	userMananger                 *UserManager
)

type UserManager struct {
	oidClients map[string]*OpenIDClient
}

func GetUserManager() *UserManager {
	return userMananger
}

func InitUserManager() {
	c := config.GetConfig()
	manager := new(UserManager)
	manager.oidClients = make(map[string]*OpenIDClient)

	for providerID := range c.GetStringMap("auth") {
		err := manager.LoadOIDProvider(providerID)

		if err != nil {
			panic(err)
		}
	}

	userMananger = manager
}

func (manager *UserManager) LoadOIDProvider(providerID string) error {
	if _, ok := manager.oidClients[providerID]; ok {
		return ErrProviderAlreadyRegistered
	}

	client, err := LoadOpenIDClient(providerID)

	if err != nil {
		return err
	}

	manager.oidClients[providerID] = client
	return nil
}

func (manager *UserManager) GetOIDClient(providerID string) (*OpenIDClient, error) {
	var client *OpenIDClient
	var ok bool

	if client, ok = manager.oidClients[providerID]; !ok {
		return nil, ErrClientNotFound
	}

	return client, nil
}

func (manager *UserManager) GenerateLoginUrl(providerID string, state string) (string, error) {
	client, err := manager.GetOIDClient(providerID)
	if err != nil {
		return "", err
	}

	return client.GenerateLoginURL(state), nil
}

func (manager *UserManager) GetAccount(providerID string, code string) (*models.UserAccount, error) {
	client, err := manager.GetOIDClient(providerID)
	if err != nil {
		return nil, err
	}

	oauthToken, err := client.FetchOAuthToken(code)
	if err != nil {
		return nil, err
	}

	oidToken, err := client.FetchOIDToken(oauthToken)
	if err != nil {
		return nil, err
	}

	account, err := manager.FindAccount(providerID, oidToken)
	if err != mgo.ErrNotFound {
		return account, err
	}

	return manager.RegisterAccount(providerID, oauthToken, oidToken)
}

func (manager *UserManager) RegisterAccount(providerID string, oauth2Token *oauth2.Token, oidToken *oidc.IDToken) (*models.UserAccount, error) {
	client, err := manager.GetOIDClient(providerID)
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

	_, err = models.RegisterOIDAccount(providerID, oidToken.Subject, account.ID)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (UserManager) FindAccount(providerID string, oidToken *oidc.IDToken) (*models.UserAccount, error) {
	oidAccount, err := models.GetOIDAccount(providerID, oidToken.Subject)
	if err != nil {
		return nil, err
	}

	account, err := models.GetUserAccount(oidAccount.Owner)
	if err != nil {
		return nil, err
	}

	return account, nil
}
