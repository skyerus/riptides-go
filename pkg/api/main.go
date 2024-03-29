package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	Router *mux.Router
	AuthRouter *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.Router.Use(Cors)
	a.AuthRouter = a.Router.PathPrefix("/api/auth").Subrouter()
	a.AuthRouter.Use(Auth)
	a.setRouters()
}

func (a *App) setRouters() {
	a.Router.HandleFunc("/", HealthCheck).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/api", ApiHealthCheck).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/api/user", CreateUser).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/api/login", Login).Methods("POST", "OPTIONS", "OPTIONS")
	a.Router.HandleFunc("/api/spotify/authorize", RedirectSpotifyAuthorize).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/api/tides", GetTides).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/api/user/{username}/forgot/password", ForgotPassword).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/api/user/reset/password", ResetPassword).Methods("POST", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/following", GetFollowing).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/followers", GetFollowers).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/following/count", GetFollowingCount).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/followers/count", GetFollowersCount).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/follow/{username}", Follow).Methods("PUT", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/unfollow/{username}", Unfollow).Methods("DELETE", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}", GetUser).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/me/config", GetMyConfig).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/spotify/authorize", AuthorizeSpotify).Methods("POST", "OPTIONS")
	a.AuthRouter.HandleFunc("/spotify/v1/me/player/play", Play).Methods("PUT", "OPTIONS")
	a.AuthRouter.HandleFunc("/spotify/v1/search", Search).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/tides", CreateTide).Methods("POST", "OPTIONS")
	a.AuthRouter.HandleFunc("/tides/genres", GetGenres).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/tides/{tideId}/favorite", FavoriteTide).Methods("PUT", "OPTIONS")
	a.AuthRouter.HandleFunc("/tides/{tideId}/unfavorite", UnfavoriteTide).Methods("PUT", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/tides/favorite", GetFavoriteTides).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/tides/favorite/count", GetFavoriteTidesCount).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/tides", GetUserTides).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/user/{username}/tides/count", GetUserTidesCount).Methods("GET", "OPTIONS")
	a.AuthRouter.HandleFunc("/avatar", UploadAvatar).Methods("POST", "OPTIONS")
	a.AuthRouter.HandleFunc("/tides", GetTidesAuth).Methods("GET", "OPTIONS")
}

func (a *App) Run(host string) {
	srv := &http.Server{
		Handler: a.Router,
		Addr: host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func openDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME"))
	if err != nil {
		log.Println(err)
		return db, err
	}
	return db, nil
}
