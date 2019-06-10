package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/skyerus/riptides-go/pkg/crypto"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type userService struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) user.Service {
	return &userService{userRepo}
}

func (u userService) Create(user models.User) customError.Error {
	hash := crypto.NewHash()
	var exists map[string]bool
	c := make(chan string, 1)
	e := make(chan error, 1)
	m := make(chan map[string]bool, 1)

	go hash.Generate(user.Password, c, e)
	go u.DoesUserExist(user, m, e)

	select {
	case user.Password = <-c:
		exists = <-m
	case err := <-e:
		return customError.NewGenericHttpError(err)
	case exists = <-m:
		user.Password = <-c
	}
	if exists["username"] {
		return customError.NewHttpError(http.StatusConflict, "A user already exists with this username", nil)
	}
	if exists["email"] {
		return customError.NewHttpError(http.StatusConflict, "A user already exists with this email", nil)
	}

	return u.userRepo.Create(user)
}

func (u userService) DoesUserExist(user models.User, m chan map[string]bool, e chan error) {
	existsMap := make(map[string]bool, 2)
	existsMap["username"] = false
	existsMap["email"] = false

	exists, err := u.userRepo.DoesUserExistWithUsername(user.Username)
	if err != nil {
		e <- err
		return
	}
	if exists {
		existsMap["username"] = true
		m <- existsMap
		return
	}

	exists, err = u.userRepo.DoesUserExistWithEmail(user.Email)
	if err != nil {
		e <- err
		return
	}
	if exists {
		existsMap["email"] = true
	}

	m <- existsMap
}

func (u userService) Authenticate(creds models.Credentials) bool {
	hash := crypto.NewHash()

	authUser := models.User{}
	authUser.Username = creds.Username
	u.userRepo.Get(&authUser)

	err := hash.Compare(authUser.Password, creds.Password)
	if err != nil {
		return false
	}
	return true
}

func (u userService) VerifyToken(token string) (*jwt.Token, error) {
	claims := &models.Claims{}

	jwtFile, err := os.Open(os.Getenv("JWT_PATH") + "/private.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer jwtFile.Close()

	jwtKey, err := ioutil.ReadAll(jwtFile)

	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}
