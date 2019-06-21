package tide

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Repository interface {
	CreateTide(user *models.User, tide *models.Tide) customError.Error
	CreateTideGenre(tide *models.Tide, genre *models.Genre) customError.Error
	GetTag(tag *models.Tag) bool
	CreateTag(tag *models.Tag) customError.Error
	CreateTideTag(tide *models.Tide, tag *models.Tag) customError.Error
	GetGenres() ([]models.Genre, customError.Error)
	GetTides(orderBy string, offset int, limit int) ([]models.Tide, customError.Error)
	GetTideParticipants(tide *models.Tide, limit int, offset int) ([]models.User, customError.Error)
	GetTideGenres(tide *models.Tide, limit int, offset int) ([]models.Genre, customError.Error)
	GetTideTags(tide *models.Tide, limit int, offset int) ([]models.Tag, customError.Error)
}
