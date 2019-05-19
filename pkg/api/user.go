package api

import (
	"database/sql"
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user/repository"
	"github.com/skyerus/riptides-go/pkg/user/service"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func CreateUser(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondBadRequest(w)
		return
	}
	if user.Username == "" || user.Email == "" || user.Password == "" {
		respondBadRequest(w)
		return
	}

	db, err := sql.Open("mysql", os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME"))
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()
	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)
	customError := userService.Create(user)

	if customError != nil {
		handleError(w, customError)
		return
	}

	respondJSON(w, 200, nil)
}

func Login(w http.ResponseWriter, r *http.Request)  {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		respondBadRequest(w)
		return
	}

}
