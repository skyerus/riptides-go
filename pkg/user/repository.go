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
	GetFollowing(user *models.User, offset int, limit int) ([]models.Following, customError.Error)
	DoesUserFollow(currentUser *models.User, user *models.User) (bool, error)
	GetFollowers(user *models.User, offset int, limit int) ([]models.Following, customError.Error)
}
