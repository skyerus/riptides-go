package service

import (
	"github.com/skyerus/riptides-go/pkg/crypto"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user"
)

type userService struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) user.Service {
	return &userService{userRepo}
}

func (u userService) Create(user models.User) customError.Error {
	hash := crypto.NewHash()

	c := make(chan string, 1)
	e := make(chan error, 1)
	m := make(chan map[string]bool, 1)
	go hash.Generate(user.Password, c, e)
	go u.DoesUserExist(user, m, e)
	return u.userRepo.Create(user, c, m, e)
}

func (u userService) DoesUserExist(user models.User, m chan map[string]bool, e chan error) {
	u.userRepo.DoesUserExist(user, m, e)
}
