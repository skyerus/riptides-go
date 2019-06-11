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
	customErr := u.userRepo.Get(&authUser)
	if customErr != nil {
		return false
	}

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

func (u userService) Get(username string) (models.User, customError.Error) {
	User := models.User{}
	User.Username = username
	customErr := u.userRepo.Get(&User)
	if customErr != nil {
		return User, customErr
	}
	return User, nil
}

func (u userService) GetCurrentUser(r *http.Request) (models.User, customError.Error) {
	User := models.User{}
	token := r.Header.Get("Authorization")

	if len(token) < 7 {
		return User, customError.NewGenericHttpError(nil)
	}
	token = token[7:]
	claims := &models.Claims{}
	Token, _, err := new(jwt.Parser).ParseUnverified(token, claims)
	if err != nil {
		return User, customError.NewGenericHttpError(nil)
	}

	tokenClaims := Token.Claims
	tokenClaims, ok := tokenClaims.(*models.Claims)
	if !ok {
		return User, customError.NewGenericHttpError(nil)
	}

	return u.Get(claims.Username)
}

func (u userService) GetMyFollowing(currentUser models.User, offset int, limit int) ([]models.Following, customError.Error) {
	return u.userRepo.GetFollowing(&currentUser, offset, limit)
}
