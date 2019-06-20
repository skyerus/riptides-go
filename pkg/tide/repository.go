package tide

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Repository interface {
	CreateTide(user *models.User, tide *models.Tide) customError.Error
	CreateTideGenre(tide *models.Tide, genre *models.Genre) customError.Error
	GetTag(name string) (models.Tag, bool)
	CreateTag(tag *models.Tag) customError.Error
	CreateTideTag(tide *models.Tide, tag *models.Tag) customError.Error
}
