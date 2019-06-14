package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user/repository"
	"github.com/skyerus/riptides-go/pkg/user/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

	db, err := openDb()
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

	db, err := openDb()
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

func GetFollowing(w http.ResponseWriter, r *http.Request) {
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
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}
	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	if CurrentUser.ID == User.ID {
		following, customErr := userService.GetMyFollowing(CurrentUser, offset, limit)
		if customErr != nil {
			respondGenericError(w)
			return
		}

		respondJSON(w, http.StatusOK, following)
		return
	}

	following, customErr := userService.GetFollowing(CurrentUser, User, offset, limit)
	if customErr != nil {
		respondGenericError(w)
		return
	}

	respondJSON(w, http.StatusOK, following)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
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
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}
	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	if CurrentUser.ID == User.ID {
		following, customErr := userService.GetMyFollowers(CurrentUser, offset, limit)
		if customErr != nil {
			respondGenericError(w)
			return
		}

		respondJSON(w, http.StatusOK, following)
		return
	}

	following, customErr := userService.GetFollowers(CurrentUser, User, offset, limit)
	if customErr != nil {
		respondGenericError(w)
		return
	}

	respondJSON(w, http.StatusOK, following)
}

func GetFollowingCount(w http.ResponseWriter, r *http.Request)  {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	followCount := models.FollowCount{}
	followCount.Count, customErr = userService.GetFollowingCount(User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, followCount)
}

func GetFollowersCount(w http.ResponseWriter, r *http.Request)  {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	followCount := models.FollowCount{}
	followCount.Count, customErr = userService.GetFollowerCount(User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, followCount)
}

func Follow(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	if CurrentUser.ID == User.ID {
		respondJSON(w, http.StatusBadRequest, nil)
		return
	}

	customErr = userService.Follow(CurrentUser, User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func Unfollow(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	if CurrentUser.ID == User.ID {
		respondJSON(w, http.StatusBadRequest, nil)
		return
	}

	customErr = userService.Unfollow(CurrentUser, User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	userService := service.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	follow := models.Following{}
	follow.User = User
	follow.Following, customErr = userService.DoesUserFollow(&CurrentUser, &User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, follow)
}
