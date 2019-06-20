package tide

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Service interface {
	CreateTide(user *models.User, tide *models.Tide) customError.Error
}