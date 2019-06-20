package SpotifyService

import (
	"encoding/base64"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/handler"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyHandler"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const SpotifyBaseUrl = "https://api.spotify.com/v1"
const SpotifyAuthorizeUrl = "https://accounts.spotify.com/api/token"

type spotifyService struct {
	spotifyRepo spotify.Repository
}

func NewSpotifyService(spotifyRepo spotify.Repository) spotify.Service {
	return &spotifyService{spotifyRepo}
}

func (s spotifyService) CredentialsExist(user *models.User) (bool, customError.Error) {
	return s.spotifyRepo.CredentialsExist(user)
}

func (s spotifyService) AuthorizeUser(user *models.User, authorization models.SpotifyAuthorization) customError.Error {
	spotifyHandler := SpotifyHandler.NewSpotifyHandler(s.spotifyRepo)
	Handler := handler.NewRequestHandler(spotifyHandler)

	body := url.Values{}
	body.Add("grant_type", "authorization_code")
	body.Add("code", authorization.Code)
	body.Add("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))

	request, err := http.NewRequest("POST", SpotifyAuthorizeUrl, strings.NewReader(body.Encode()))
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	clientData := os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET")
	clientBase64 := base64.StdEncoding.EncodeToString([]byte(clientData))

	request.Header.Add("Authorization", "Basic " + clientBase64)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, customErr := Handler.SendRequest(request, user, false, false)
	if customErr != nil {
		return customErr
	}

	return spotifyHandler.SaveCredentials(response, user)
}
