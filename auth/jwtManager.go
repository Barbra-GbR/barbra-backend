package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/Barbra-GbR/barbra-backend/config"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrSecretInvalid           = errors.New("jwt: len of secret should be > 10")
	ErrUnexpectedSigningMethod = errors.New("jwt: unexpected signing method")
	ErrTokenInvalid            = errors.New("jwt: invalid token")
	ErrTokenExpired            = errors.New("jwt: token expired")
	jwtManager                 *JWTManager
)

//Provides tools to validate jwt tokens and retrieve and getting corresponding accounts
type JWTManager struct {
	secret  []byte
	issuer  string
	idCount int64
}

//Returns the initialized JWT. Do not call before calling InitializeAccountManager!
func GetJWT() *JWTManager {
	return jwtManager
}

//Initialises the JWTManager with data from the config
func InitializeJWT() {
	c := config.GetConfig()

	j := new(JWTManager)
	j.secret = []byte(c.GetString("server.jwt_secret"))
	j.issuer = c.GetString("server.host")
	j.idCount = 0
	if len(j.secret) < 10 {
		panic(ErrSecretInvalid)
	}

	jwtManager = j
}

//Generates a new JWT for the given information
func (manager *JWTManager) GenerateToken(accountId bson.ObjectId) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId.Hex(),
		Issuer:    manager.issuer,
		ExpiresAt: now.Add(time.Hour * 24 * 30).Unix(),
		IssuedAt:  now.Unix(),
		Audience:  manager.issuer,
		NotBefore: now.Unix(),
		Id:        manager.NewTokenId(accountId),
	})

	return token.SignedString(manager.secret)
}

//Validates the token and returns the corresponding userAccount
func (manager *JWTManager) GetAccountId(tokenString string) (string, error) {
	claims := new(jwt.StandardClaims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return manager.secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrTokenInvalid
	}

	now := time.Now().Unix()

	//TODO: Id disabling of tokens (maybe)
	if claims.NotBefore >= now || claims.Audience != manager.issuer {
		return "", ErrTokenInvalid
	}
	if claims.ExpiresAt <= now {
		return "", ErrTokenExpired
	}

	return claims.Subject, nil
}

//Generates a new TokenId for use within JWTs
func (manager *JWTManager) NewTokenId(accountId bson.ObjectId) string {
	m := md5.New()
	m.Write([]byte(accountId.Hex() + ":" + time.Now().String() + ":" + string(manager.idCount)))
	return hex.EncodeToString(m.Sum(nil))
}
