package spotify

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Repository interface {
	CredentialsExist(user *models.User) (bool, customError.Error)
}