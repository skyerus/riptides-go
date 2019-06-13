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
	a.AuthRouter = a.Router.PathPrefix("/api/auth").Subrouter()
	a.AuthRouter.Use(Auth)
	a.setRouters()
}

func (a *App) setRouters() {
	a.Router.HandleFunc("/healthcheck", HealthCheck).Methods("GET")
	a.Router.HandleFunc("/api/user", CreateUser).Methods("POST")
	a.Router.HandleFunc("/api/login", Login).Methods("POST")
	a.AuthRouter.HandleFunc("/user/{username}/following", GetFollowing).Methods("GET")
	a.AuthRouter.HandleFunc("/user/{username}/followers", GetFollowers).Methods("GET")
	a.AuthRouter.HandleFunc("/user/{username}/following/count", GetFollowingCount).Methods("GET")
	a.AuthRouter.HandleFunc("/user/{username}/followers/count", GetFollowersCount).Methods("GET")
	a.AuthRouter.HandleFunc("/user/follow/{username}", Follow).Methods("PUT")
	a.AuthRouter.HandleFunc("/user/unfollow/{username}", Unfollow).Methods("DELETE")
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
	return sql.Open("mysql", os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME"))
}
