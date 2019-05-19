package user

import "github.com/skyerus/riptides-go/pkg/models"

type Service interface {
	Create(user models.User) error
}