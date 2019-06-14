package api

import (
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
