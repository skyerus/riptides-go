package api

import (
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyRepository"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyService"
	"github.com/skyerus/riptides-go/pkg/user/UserRepository"
	"github.com/skyerus/riptides-go/pkg/user/UserService"
	"log"
	"net/http"
	"net/url"
	"os"
)

const AuthorizeUrl = "https://accounts.spotify.com/authorize"

func RedirectSpotifyAuthorize(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", AuthorizeUrl, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	query := url.Values{}
	query.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	query.Add("response_type", "code")
	query.Add("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))
	query.Add("scope", "user-modify-playback-state")
	req.URL.RawQuery = query.Encode()

	http.Redirect(w, r, req.URL.String(), http.StatusSeeOther)
}

func AuthorizeSpotify(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var auth models.SpotifyAuthorization
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		respondBadRequest(w)
		return
	}

	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	spotifyRepo := SpotifyRepository.NewMysqlSpotifyRepository(db)
	userService := UserService.NewUserService(userRepo)
	spotifyService := SpotifyService.NewSpotifyService(spotifyRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	customErr = spotifyService.AuthorizeUser(&CurrentUser, auth)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func Play(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload models.Play
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		respondBadRequest(w)
		return
	}

	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	spotifyRepo := SpotifyRepository.NewMysqlSpotifyRepository(db)
	userService := UserService.NewUserService(userRepo)
	spotifyService := SpotifyService.NewSpotifyService(spotifyRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	var spotifyPlay models.SpotifyPlay
	var URIs [1]string
	URIs[0] = payload.URI
	spotifyPlay.URIs = URIs

	customErr = spotifyService.Play(&CurrentUser, spotifyPlay)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func Search(w http.ResponseWriter, r *http.Request)  {
	rawQuery := r.URL.RawQuery

	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	spotifyRepo := SpotifyRepository.NewMysqlSpotifyRepository(db)
	userService := UserService.NewUserService(userRepo)
	spotifyService := SpotifyService.NewSpotifyService(spotifyRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	simpleSpotifySearch, customErr := spotifyService.Search(&CurrentUser, rawQuery)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, simpleSpotifySearch)
}
