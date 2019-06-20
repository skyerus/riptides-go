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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const SpotifyBaseUrl = "https://api.spotify.com/v1"
const SpotifyAuthorizeUrl = "https://accounts.spotify.com/api/token"
const SpotifyPlayEndpoint = "/me/player/play"
const SpotifySearchEndpoint = "/search"

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

func (s spotifyService) Play(user *models.User, spotifyPlay models.SpotifyPlay) customError.Error {
	spotifyHandler := SpotifyHandler.NewSpotifyHandler(s.spotifyRepo)
	Handler := handler.NewRequestHandler(spotifyHandler)

	bodyBytes, err := json.Marshal(spotifyPlay)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	b := bytes.NewBuffer(bodyBytes)

	request, err := http.NewRequest("PUT", SpotifyBaseUrl + SpotifyPlayEndpoint, b)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	_, customErr := Handler.SendRequest(request, user, true, true)

	return customErr
}

func (s spotifyService) Search(user *models.User, query string) (models.SpotifySearchSimple, customError.Error) {
	var simpleSearchResponse models.SpotifySearchSimple
	spotifyHandler := SpotifyHandler.NewSpotifyHandler(s.spotifyRepo)
	Handler := handler.NewRequestHandler(spotifyHandler)

	request, err := http.NewRequest("GET", SpotifyBaseUrl + SpotifySearchEndpoint + "?" + query, nil)
	if err != nil {
		return simpleSearchResponse, customError.NewGenericHttpError(err)
	}

	response, customErr := Handler.SendRequest(request, user, true, true)
	if customErr != nil {
		return simpleSearchResponse, customErr
	}
	defer response.Body.Close()

	var spotifySearch models.SpotifySearch
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return simpleSearchResponse, customError.NewGenericHttpError(err)
	}
	err = json.Unmarshal(body, &spotifySearch)
	if err != nil {
		return simpleSearchResponse, customError.NewGenericHttpError(err)
	}

	if len(spotifySearch.Tracks.Items) < 1 {
		return simpleSearchResponse, customError.NewHttpError(http.StatusNotFound, "No song found", nil)
	}
	simpleSearchResponse.URI = spotifySearch.Tracks.Items[0].URI
	simpleSearchResponse.Artist = spotifySearch.Tracks.Items[0].Artists[0].Name
	simpleSearchResponse.Name = spotifySearch.Tracks.Items[0].Name
	simpleSearchResponse.DurationMs = spotifySearch.Tracks.Items[0].DurationMs
	if len(spotifySearch.Tracks.Items[0].Album.Images) < 3 {
		return simpleSearchResponse, customError.NewHttpError(http.StatusNotFound, "No album cover found", nil)
	}
	simpleSearchResponse.Image = spotifySearch.Tracks.Items[0].Album.Images[2]

	return simpleSearchResponse, nil
}
