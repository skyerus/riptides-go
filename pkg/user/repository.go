package user

import "github.com/skyerus/riptides-go/pkg/models"

type Repository interface {
	Create(user models.User) error
	DoesUserExist(user models.User) error
}
