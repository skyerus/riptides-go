package api

import (
	"database/sql"
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user/repository"
	"github.com/skyerus/riptides-go/pkg/user/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func CreateUser(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println(err)
		respondGenericError(w)
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
	err = userService.Create(user)

	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}

	respondJSON(w, 200, nil)
}