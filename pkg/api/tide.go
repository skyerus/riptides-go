package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/notifications"
	"github.com/skyerus/riptides-go/pkg/tide/TideRepository"
	"github.com/skyerus/riptides-go/pkg/tide/TideService"
	"github.com/skyerus/riptides-go/pkg/user/UserRepository"
	"github.com/skyerus/riptides-go/pkg/user/UserService"
	"net/http"
	"strconv"
)

func CreateTide(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	var Tide models.Tide

	err := json.NewDecoder(r.Body).Decode(&Tide)
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
	userService := UserService.NewUserService(userRepo)
	tideRepo := TideRepository.NewMysqlTideRepository(db)
	tideService := TideService.NewTideService(tideRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	customErr = tideService.CreateTide(&CurrentUser, &Tide)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func GetTides(w http.ResponseWriter, r *http.Request)  {
	params := r.URL.Query()
	offset, err := strconv.Atoi(params.Get("offset"))
	if err != nil {
		respondBadRequest(w)
		return
	}
	limit, err := strconv.Atoi(params.Get("limit"))
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

	tideRepo := TideRepository.NewMysqlTideRepository(db)
	tideService := TideService.NewTideService(tideRepo)

	tides, customErr := tideService.GetTides("date_created", offset, limit)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, tides)
}

func GetGenres(w http.ResponseWriter, r *http.Request) {
	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	tideRepo := TideRepository.NewMysqlTideRepository(db)
	tideService := TideService.NewTideService(tideRepo)

	genres, customErr := tideService.GetGenres()
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, genres)
}

func FavoriteTide(w http.ResponseWriter, r *http.Request) {
	tideId, err := strconv.Atoi(mux.Vars(r)["tideId"])
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
	userService := UserService.NewUserService(userRepo)
	tideRepo := TideRepository.NewMysqlTideRepository(db)
	tideService := TideService.NewTideService(tideRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	Tide, customErr := tideService.GetTide(tideId)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	exists, customErr := tideService.IsTideFavorited(&Tide, &CurrentUser)
	if customErr != nil {
		handleError(w, customErr)
		return
	}
	if !exists {
		customErr = tideService.FavoriteTide(&Tide, &CurrentUser)
		if customErr != nil {
			handleError(w, customErr)
			return
		}

		User, customErr := userService.GetFromId(Tide.User.ID)
		if customErr != nil {
			handleError(w, customErr)
			return
		}

		token, customErr := userService.GenerateToken(User.Username)
		if customErr != nil {
			handleError(w, customErr)
			return
		}
		go notifications.PushNotification(token, CurrentUser.Username + " favorited your tide")
	}

	respondJSON(w, http.StatusOK, nil)
}

func UnfavoriteTide(w http.ResponseWriter, r *http.Request)  {
	tideId, err := strconv.Atoi(mux.Vars(r)["tideId"])
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
	userService := UserService.NewUserService(userRepo)
	tideRepo := TideRepository.NewMysqlTideRepository(db)
	tideService := TideService.NewTideService(tideRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	Tide, customErr := tideService.GetTide(tideId)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	customErr = tideService.UnfavoriteTide(&Tide, &CurrentUser)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func GetFavoriteTides(w http.ResponseWriter, r *http.Request) {
	username, success := mux.Vars(r)["username"]
	if !success {
		respondBadRequest(w)
		return
	}

	params := r.URL.Query()
	offset, err := strconv.Atoi(params.Get("offset"))
	if err != nil {
		respondBadRequest(w)
		return
	}
	limit, err := strconv.Atoi(params.Get("limit"))
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
	userService := UserService.NewUserService(userRepo)
	tideRepo := TideRepository.NewMysqlTideRepository(db)
	tideService := TideService.NewTideService(tideRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	tides, customErr := tideService.GetFavoriteTides(&User, offset, limit)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, tides)
}

func GetFavoriteTidesCount(w http.ResponseWriter, r *http.Request)  {
	username, success := mux.Vars(r)["username"]
	if !success {
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
	userService := UserService.NewUserService(userRepo)
	tideRepo := TideRepository.NewMysqlTideRepository(db)
	tideService := TideService.NewTideService(tideRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	count := models.Count{}
	count.Count, customErr = tideService.GetFavoriteTidesCount(&User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, count)
}
