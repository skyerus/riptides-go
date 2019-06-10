package api

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user/repository"
	"github.com/skyerus/riptides-go/pkg/user/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

	respondJSON(w, http.StatusOK, nil)
}

func Login(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
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

	if !userService.Authenticate(creds) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtFile, err := os.Open(os.Getenv("JWT_PATH") + "/private.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer jwtFile.Close()

	jwtKey, err := ioutil.ReadAll(jwtFile)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make(map[string]string, 1)
	response["token"] = tokenString

	respondJSON(w, http.StatusOK, response)
}
