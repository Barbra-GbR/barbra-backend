package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/bitphinix/babra_backend/config"
)

var (
	ErrSecretInvalid           = errors.New("jwt: len of secret should be > 10")
	ErrUnexpectedSigningMethod = errors.New("jwt: unexpected signing method")
	ErrTokenInvalid            = errors.New("jwt: invalid token")
	ErrTokenExpired            = errors.New("jwt: token expired")
	jwtHandler                 *JWT
)

type JWT struct {
	secret  []byte
	issuer  string
	idCount int64
}

func GetJWT() *JWT {
	return jwtHandler
}

func InitJWT() {
	c := config.GetConfig()

	j := new(JWT)
	j.secret = []byte(c.GetString("server.jwt_secret"))
	j.issuer = c.GetString("server.host")
	j.idCount = 0

	if len(j.secret) < 10 {
		panic(ErrSecretInvalid)
	}

	jwtHandler = j
}

func (j *JWT) GenerateToken(accountId string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId,
		Issuer:    j.issuer,
		ExpiresAt: now.Add(time.Hour * 24 * 30).Unix(),
		IssuedAt:  now.Unix(),
		Audience:  j.issuer,
		NotBefore: now.Unix(),
		Id:        j.NewTokenID(accountId),
	})

	return token.SignedString(j.secret)
}

func (j *JWT) GetUserId(tokenString string) (string, error) {
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return j.secret, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrTokenInvalid
	}

	now := time.Now().Unix()

	//TODO: Id disabling of tokens (maybe)
	if claims.NotBefore >= now || claims.Audience != j.issuer {
		return "", ErrTokenInvalid
	}

	if claims.ExpiresAt <= now {
		return "", ErrTokenExpired
	}

	return claims.Subject, nil
}

func (j *JWT) NewTokenID(accountId string) string {
	m := md5.New()
	m.Write([]byte(accountId + ":" + time.Now().String() + ":" + string(j.idCount)))
	return hex.EncodeToString(m.Sum(nil))
}
