package SpotifyService

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/handler"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyHandler"
	"net/http"
	"os"
)

const SpotifyBaseUrl = "https://api.spotify.com/v1"
const SpotifyAuthorizeUrl = "https://accounts.spotify.com/api/token"

type authorizationBody struct {
	GrantType string `json:"grant_type"`
	Code string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

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
	authBody := authorizationBody{"authorization_code", authorization.Code, os.Getenv("SPOTIFY_REDIRECT_URI")}
	spotifyHandler := SpotifyHandler.NewSpotifyHandler(s.spotifyRepo)
	Handler := handler.NewRequestHandler(spotifyHandler)

	bodyBytes, err := json.Marshal(authBody)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	b := bytes.NewBuffer(bodyBytes)
	request, err := http.NewRequest("POST", SpotifyAuthorizeUrl, b)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	clientData := os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET")
	clientBase64 := base64.StdEncoding.EncodeToString([]byte(clientData))

	request.Header.Set("Authorization", "Basic " + clientBase64)

	response, customErr := Handler.SendRequest(request, user, false, false)
	if customErr != nil {
		return customErr
	}

	return spotifyHandler.SaveCredentials(response, user)
}
