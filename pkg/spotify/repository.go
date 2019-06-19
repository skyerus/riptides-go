package spotify

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyHandler"
)

type Repository interface {
	CredentialsExist(user *models.User) (bool, customError.Error)
	GetCredentials(user *models.User) (models.SpotifyCredentials, customError.Error)
	SaveCredentials(creds SpotifyHandler.Credentials, user *models.User) customError.Error
}