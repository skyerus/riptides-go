package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	Salt string `json:"-"`
	Avatar string `json:"avatar"`
	Bio string `json:"bio"`
}

type Following struct {
	User User `json:"user"`
	Following bool `json:"following"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Config struct {
	Spotify bool `json:"spotify"`
}

type UserConfig struct {
	User `json:"user"`
	Config `json:"config"`
}

type ResetPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token string `json:"token"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	return c.StandardClaims.Valid()
}

