package SpotifyService

import (
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify"
)

type spotifyService struct {
	spotifyRepo spotify.Repository
}

func NewSpotifyService(spotifyRepo spotify.Repository) spotify.Service {
	return &spotifyService{spotifyRepo}
}

func (s spotifyService) CredentialsExist(user *models.User) (bool, customError.Error) {
	return s.spotifyRepo.CredentialsExist(user)
}
