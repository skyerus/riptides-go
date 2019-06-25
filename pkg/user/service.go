package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"net/http"
)

type Service interface {
	Create(user models.User) customError.Error
	DoesUserExist(user models.User, m chan map[string]bool, e chan error)
	Authenticate(creds models.Credentials) bool
	VerifyToken(token string) (*jwt.Token, error)
	Get(username string) (models.User, customError.Error)
	GetFromId(id int) (models.User, customError.Error)
	GetCurrentUser(r *http.Request) (models.User, customError.Error)
	GetMyFollowing(currentUser models.User, offset int, limit int) ([]models.Following, customError.Error)
	GetFollowing(currentUser models.User, user models.User, offset int, limit int) ([]models.Following, customError.Error)
	GetMyFollowers(currentUser models.User, offset int, limit int) ([]models.Following, customError.Error)
	GetFollowers(currentUser models.User, user models.User, offset int, limit int) ([]models.Following, customError.Error)
	GetFollowingCount(user models.User) (int, customError.Error)
	GetFollowerCount(user models.User) (int, customError.Error)
	Follow(currentUser models.User, user models.User) customError.Error
	Unfollow(currentUser models.User, user models.User) customError.Error
	DoesUserFollow(currentUser *models.User, user *models.User) (bool, customError.Error)
	GenerateToken(username string) (string, customError.Error)
	SaveAvatar(currentUser *models.User) customError.Error
}