package spotify

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
)

type Service interface {
	CredentialsExist(user *models.User) (bool, customError.Error)
	AuthorizeUser(user *models.User, authorization models.SpotifyAuthorization) customError.Error
	Play(user *models.User, spotifyPlay models.SpotifyPlay) customError.Error
	Search(user *models.User, query string) (models.SpotifySearchSimple, customError.Error)
}
