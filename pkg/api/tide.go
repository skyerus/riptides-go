package api

import (
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/tide/TideRepository"
	"github.com/skyerus/riptides-go/pkg/tide/TideService"
	"github.com/skyerus/riptides-go/pkg/user/UserRepository"
	"github.com/skyerus/riptides-go/pkg/user/UserService"
	"net/http"
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

}
