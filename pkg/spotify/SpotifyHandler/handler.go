package SpotifyHandler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/handler"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type spotifyHandler struct {
	spotifyRepo spotify.Repository
}

func NewSpotifyHandler(spotifyRepo spotify.Repository) handler.Handler {
	return &spotifyHandler{spotifyRepo}
}

func (handler spotifyHandler) SaveCredentials(response *http.Response, user *models.User) customError.Error {
	var creds models.SpotifyCredentials
	err := json.NewDecoder(response.Body).Decode(&creds)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	credsExist, customErr := handler.spotifyRepo.CredentialsExist(user)
	if customErr != nil {
		return customErr
	}

	if credsExist {
		return handler.spotifyRepo.UpdateCredentials(creds, user)
	}
	return handler.spotifyRepo.SaveCredentials(creds, user)
}

func (handler spotifyHandler) HandleAuthorizedRequest(r *http.Request, user *models.User) customError.Error {
	creds, customErr := handler.spotifyRepo.GetCredentials(user)
	if customErr != nil {
		return customErr
	}
	r.Header.Set("Authorization", "Bearer " + creds.AccessToken)

	return nil
}

func (handler spotifyHandler) GetRefreshRequest(user *models.User) (*http.Request, customError.Error) {
	creds, customErr := handler.spotifyRepo.GetCredentials(user)
	if customErr != nil {
		return nil, customErr
	}

	body := url.Values{}
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", creds.RefreshToken)

	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, customError.NewGenericHttpError(err)
	}

	clientData := os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET")
	clientBase64 := base64.StdEncoding.EncodeToString([]byte(clientData))

	request.Header.Add("Authorization", "Basic " + clientBase64)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

