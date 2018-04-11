package auth

import (
	"errors"
	"gopkg.in/mgo.v2"
	"golang.org/x/oauth2"
	"github.com/coreos/go-oidc"
	"github.com/Barbra-GbR/barbra-backend/config"
	"github.com/Barbra-GbR/barbra-backend/models"
)

var (
	ErrProviderAlreadyRegistered = errors.New("provider is already registered")
	ErrClientNotFound            = errors.New("no client with the specified providerId registered")
	accountManager               *OpenIdManager
)

//Provides tools for registering and finding users with OpenId, OAuth2 tokens
type OpenIdManager struct {
	oidClients map[string]*OpenIdClient
}

//Returns the initialized OpenIdManager. Do not call before calling InitializeAccountManager!
func GetAccountManager() *OpenIdManager {
	return accountManager
}

//Initialises the OpenIdManager with data from the config
func InitializeAccountManager() {
	c := config.GetConfig()
	manager := new(OpenIdManager)
	manager.oidClients = make(map[string]*OpenIdClient)

	for providerId := range c.GetStringMap("auth") {
		err := manager.LoadOIdProvider(providerId)
		if err != nil {
			panic(err)
		}
	}

	accountManager = manager
}

//Loads the OIdProvider with the corresponding providerId out of the config
func (manager *OpenIdManager) LoadOIdProvider(providerId string) error {
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

//Returns the OIdClient with the given providerId
func (manager *OpenIdManager) GetOIdClient(providerId string) (*OpenIdClient, error) {
	var client *OpenIdClient
	var ok bool

	if client, ok = manager.oidClients[providerId]; !ok {
		return nil, ErrClientNotFound
	}

	return client, nil
}

//Generates a new login url with the given state
func (manager *OpenIdManager) GenerateLoginUrl(providerId string, state string) (string, error) {
	client, err := manager.GetOIdClient(providerId)
	if err != nil {
		return "", err
	}

	return client.GenerateLoginURL(state), nil
}

//Redeems the oauth code and try's to find the belonging UserAccount. If no account was found, a new one will be registered
func (manager *OpenIdManager) GetAccount(providerId string, code string) (*models.UserAccount, error) {
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

	account, err := manager.GetAccountByIdToken(providerId, oidToken)
	if err != mgo.ErrNotFound {
		return account, err
	}

	return manager.RegisterAccount(providerId, oauthToken, oidToken)
}

//Registers a new UserAccount by fetching data with the oauth2Token and the oidToken
func (manager *OpenIdManager) RegisterAccount(providerId string, oauth2Token *oauth2.Token, oidToken *oidc.IDToken) (*models.UserAccount, error) {
	client, err := manager.GetOIdClient(providerId)
	if err != nil {
		return nil, err
	}

	payload, err := client.FetchProfilePayload(oauth2Token)
	if err != nil {
		return nil, err
	}

	account, err := models.RegisterUser(payload)
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

//Try's to find a UserAccount with the give providerId and oidToken
func (manager *OpenIdManager) GetAccountByIdToken(providerId string, oidToken *oidc.IDToken) (*models.UserAccount, error) {
	oidAccount, err := models.FindOIdAccount(providerId, oidToken.Subject)
	if err != nil {
		return nil, err
	}

	account, err := models.GetUserAccountById(oidAccount.OwnerId)
	return account, err
}
