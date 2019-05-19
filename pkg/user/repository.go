package user

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Repository interface {
	Create(user models.User, c chan string, m chan map[string]bool, e chan error) customError.Error
	DoesUserExist(user models.User, m chan map[string]bool, e chan error)
}
