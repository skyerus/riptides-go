package user

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Repository interface {
	Create(user models.User) customError.Error
	DoesUserExistWithUsername(username string) (bool, error)
	DoesUserExistWithEmail(email string) (bool, error)
	Get(user *models.User) customError.Error
}
