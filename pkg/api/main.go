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
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Router.PathPrefix("/api/auth/").Subrouter().Use(Auth)
	a.Get("/healthcheck", HealthCheck)
	a.Post("/api/user", CreateUser)
	a.Post("/api/login", Login)
	a.Get("/api/auth/user/{username}/following", GetFollowing)
	a.Router.Use()
}

// HTTP wrappers
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
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
