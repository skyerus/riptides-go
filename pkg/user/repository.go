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
	GetFromId(id int) (models.User, customError.Error)
	GetFollowing(user *models.User, offset int, limit int) ([]models.Following, customError.Error)
	DoesUserFollow(currentUser *models.User, user *models.User) (bool, customError.Error)
	GetFollowers(user *models.User, offset int, limit int) ([]models.Following, customError.Error)
	GetFollowingCount(user *models.User) (int, customError.Error)
	GetFollowerCount(user *models.User) (int, customError.Error)
	Follow(currentUser *models.User, user *models.User) customError.Error
	Unfollow(currentUser *models.User, user *models.User) customError.Error
}
