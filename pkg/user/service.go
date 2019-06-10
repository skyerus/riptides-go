package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Service interface {
	Create(user models.User) customError.Error
	DoesUserExist(user models.User, m chan map[string]bool, e chan error)
	Authenticate(creds models.Credentials) bool
	VerifyToken(token string) (*jwt.Token, error)
}