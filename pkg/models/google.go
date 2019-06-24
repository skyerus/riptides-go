package models

import "github.com/dgrijalva/jwt-go"

type GoogleClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func (c GoogleClaims) Valid() error {
	return c.StandardClaims.Valid()
}

type GoogleCredentials struct {
	ClientEmail string `json:"client_email"`
	PrivateKey string `json:"private_key"`
}

type GoogleAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
}
