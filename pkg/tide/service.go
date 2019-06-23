package tide

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Service interface {
	CreateTide(user *models.User, tide *models.Tide) customError.Error
	GetGenres() ([]models.Genre, customError.Error)
	GetTides(orderBy string, offset int, limit int) ([]models.Tide, customError.Error)
	FavoriteTide(tide *models.Tide, user *models.User) customError.Error
	UnfavoriteTide(tide *models.Tide, user *models.User) customError.Error
	IsTideFavorited(tide *models.Tide, user *models.User) (bool, customError.Error)
	GetTide(id int) (models.Tide, customError.Error)
	GetFavoriteTides(user *models.User, offset int, limit int) ([]models.Tide, customError.Error)
	GetFavoriteTidesCount(user *models.User) (int, customError.Error)
	GetUserTides(user *models.User, offset int, limit int) ([]models.Tide, customError.Error)
	GetUserTidesCount(user *models.User) (*models.Count, customError.Error)
}