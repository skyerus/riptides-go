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

	for _, genreId := range tide.Genres {
		customErr = t.tideRepo.CreateTideGenre(tide, &models.Genre{ID: genreId})
		if customErr != nil {
			return customErr
		}
	}

	for _, tagName := range tide.Tags {
		tag, exists := t.tideRepo.GetTag(tagName)
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
