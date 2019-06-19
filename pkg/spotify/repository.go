package spotify

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Repository interface {
	CredentialsExist(user *models.User) (bool, customError.Error)
	GetCredentials(user *models.User) (models.SpotifyCredentials, customError.Error)
	SaveCredentials(creds models.SpotifyCredentials, user *models.User) customError.Error
	UpdateCredentials(creds models.SpotifyCredentials, user *models.User) customError.Error
}