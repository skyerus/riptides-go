package SpotifyHandler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/handler"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify"
	"net/http"
	"os"
)

type refreshBody struct {
	GrantType string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type spotifyHandler struct {
	spotifyRepo spotify.Repository
}

func NewSpotifyHandler(spotifyRepo spotify.Repository) handler.Handler {
	return &spotifyHandler{spotifyRepo}
}

func (handler spotifyHandler) SaveCredentials(response *http.Response, user *models.User) customError.Error {
	var creds models.SpotifyCredentials
	err := json.NewDecoder(response.Body).Decode(creds)
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
	var refreshBody refreshBody
	refreshBody.GrantType = "refresh_token"

	creds, customErr := handler.spotifyRepo.GetCredentials(user)
	if customErr != nil {
		return nil, customErr
	}
	refreshBody.RefreshToken = creds.RefreshToken

	bodyBytes, err := json.Marshal(refreshBody)
	if err != nil {
		return nil, customError.NewGenericHttpError(err)
	}
	b := bytes.NewBuffer(bodyBytes)
	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", b)
	if err != nil {
		return nil, customError.NewGenericHttpError(err)
	}

	clientId := []byte(os.Getenv("SPOTIFY_CLIENT_ID"))
	clientSecret := []byte(os.Getenv("SPOTIFY_CLIENT_SECRET"))
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	_, err = encoder.Write(clientId)
	if err != nil {
		return nil, customError.NewGenericHttpError(err)
	}
	err = encoder.Close()
	if err != nil {
		return nil, customError.NewGenericHttpError(err)
	}
	_, err = encoder.Write(clientSecret)
	if err != nil {
		return nil, customError.NewGenericHttpError(err)
	}
	err = encoder.Close()
	if err != nil {
		return nil, customError.NewGenericHttpError(err)
	}

	request.Header.Set("Authorization", "Basic " + string(clientId) + ":" + string(clientSecret))

	return request, nil
}

