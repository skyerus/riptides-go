package TideService

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/tide"
)

type tideService struct {
	tideRepo tide.Repository
}

func NewTideService(tideRepo tide.Repository) tide.Service {
	return &tideService{tideRepo}
}

func (t tideService) CreateTide(user *models.User, tide *models.Tide) customError.Error {
	customErr := t.tideRepo.CreateTide(user, tide)
	if customErr != nil {
		return customErr
	}

	for _, genre := range tide.Genres {
		customErr = t.tideRepo.CreateTideGenre(tide, &genre)
		if customErr != nil {
			return customErr
		}
	}

	for _, tag := range tide.Tags {
		exists := t.tideRepo.GetTag(&tag)
		if !exists {
			customErr = t.tideRepo.CreateTag(&tag)
			if customErr != nil {
				return customErr
			}
		}
		customErr = t.tideRepo.CreateTideTag(tide, &tag)
		if customErr != nil {
			return customErr
		}
	}

	return nil
}

func (t tideService) GetGenres() ([]models.Genre, customError.Error) {
	return t.tideRepo.GetGenres()
}

func (t tideService) GetTides(orderBy string, offset int, limit int) ([]models.Tide, customError.Error) {
	return t.tideRepo.GetTides(orderBy, offset, limit)
}
