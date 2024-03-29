package api

import (
	"bytes"
	"encoding/json"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/skyerus/riptides-go/pkg/RedisClient"
	"github.com/skyerus/riptides-go/pkg/crypto"
	"github.com/skyerus/riptides-go/pkg/email"
	"github.com/skyerus/riptides-go/pkg/google/GoogleHandler"
	"github.com/skyerus/riptides-go/pkg/handler"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/notifications"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyRepository"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyService"
	"github.com/skyerus/riptides-go/pkg/user/UserRepository"
	"github.com/skyerus/riptides-go/pkg/user/UserService"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
		respondGenericError(w)
		return
	}
	defer db.Close()
	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)
	user.Avatar = "https://www.synbio.cam.ac.uk/images/avatar-generic.jpg/image_preview"
	customError := userService.Create(user)

	if customError != nil {
		handleError(w, customError)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
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

	if !userService.Authenticate(creds) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, customErr := userService.GenerateToken(creds.Username)
	if customErr != nil {
		handleError(w, customErr)
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
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

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
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

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

func GetFollowingCount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	followCount := models.Count{}
	followCount.Count, customErr = userService.GetFollowingCount(User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, followCount)
}

func GetFollowersCount(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	followCount := models.Count{}
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
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

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

	doesUserFollow, customErr := userService.DoesUserFollow(&CurrentUser, &User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	if doesUserFollow {
		respondJSON(w, http.StatusOK, nil)
		return
	}

	customErr = userService.Follow(CurrentUser, User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	token, customErr := userService.GenerateToken(User.Username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}
	go notifications.PushNotification(token, CurrentUser.Username + " followed you")

	respondJSON(w, http.StatusOK, nil)
}

func Unfollow(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

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
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

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

func GetMyConfig(w http.ResponseWriter, r *http.Request) {
	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	spotifyRepo := SpotifyRepository.NewMysqlSpotifyRepository(db)
	spotifyService := SpotifyService.NewSpotifyService(spotifyRepo)

	userConfig := models.UserConfig{}
	userConfig.User = CurrentUser
	userConfig.Config.Spotify, customErr = spotifyService.CredentialsExist(&CurrentUser)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, userConfig)
}

func UploadAvatar(w http.ResponseWriter, r *http.Request)  {
	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

	CurrentUser, customErr := userService.GetCurrentUser(r)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	err = r.ParseMultipartForm(1000000)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Images are limited to 1MB"})
		return
	}

	file, fileHandler, err := r.FormFile("file")
	if err != nil {
		respondBadRequest(w)
		return
	}
	defer file.Close()

	mime := fileHandler.Header.Get("Content-Type")
	if mime != "image/png" && mime != "image/jpeg" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Please use an image of type jpeg or png"})
		return
	}

	fileName := fileHandler.Filename + strconv.Itoa(int(time.Now().Unix()))

	redisClient := RedisClient.NewRedisClient()
	defer redisClient.Close()

	googleHandler := GoogleHandler.NewGoogleHandler(redisClient)
	Handler := handler.NewRequestHandler(googleHandler)

	req, err := http.NewRequest("POST", models.GoogleStorageUploadUrl + "&name=" + fileName, file)
	if err != nil {
		respondGenericError(w)
		return
	}

	req.Header.Add("Content-Type", mime)
	req.Header.Add("Content-Length", strconv.Itoa(int(fileHandler.Size)))

	var googleResponse models.GoogleUploadResponse
	response, customErr := Handler.SendRequest(req, &CurrentUser, true, true)
	err = json.NewDecoder(response.Body).Decode(&googleResponse)
	if err != nil {
		respondGenericError(w)
		return
	}

	CurrentUser.Avatar = googleResponse.MediaLink
	customErr = userService.SaveAvatar(&CurrentUser)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	db, err := openDb()
	if err != nil {
		respondGenericError(w)
		return
	}
	defer db.Close()

	userRepo := UserRepository.NewMysqlUserRepository(db)
	userService := UserService.NewUserService(userRepo)

	User, customErr := userService.Get(username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	key := username + "_password_token"
	redisClient := RedisClient.NewRedisClient()
	defer redisClient.Close()

	_, err = redisClient.Get(key).Result()
	if err == nil {
		respondError(w, http.StatusUnprocessableEntity, "A reset email has been issued recently, please check your spam folder")
		return
	} else if err != redis.Nil {
		log.Println(err)
		respondGenericError(w)
		return
	}

	token := make([]byte, 40)
	rand.Read(token)
	hashedBytes, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	s := time.Now().Unix()
	tokenString := string(hashedBytes) + strconv.Itoa(int(s))

	redisClient.Set(key, tokenString, time.Duration(3600) * time.Second)

	link := os.Getenv("RIPTIDES_HOST") + "/reset/password?" + "token=" + tokenString + "&username=" + username

	tmpl := template.Must(template.ParseFiles(os.Getenv("FORGOT_PASSWORD_TPL")))

	var tplBuffer bytes.Buffer
	err = tmpl.Execute(&tplBuffer, models.ForgotPasswordData{
		Username: username,
		Link: link,
	})

	mg := email.NewMailgun()

	emailStr := tplBuffer.String()
	message := mg.NewMessage("riptides <noreply@riptides.io>", "Forgot password", "", User.Email)
	message.SetHtml(emailStr)

	_, _, err = mg.Send(message)
	if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]string{"message": "An email has been sent with further instructions, please check your spam folder."})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reset models.ResetPassword
	err := json.NewDecoder(r.Body).Decode(&reset)
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

	User, customErr := userService.Get(reset.Username)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	c := make(chan string, 1)
	e := make(chan error, 1)
	hash := crypto.Hash{}
	go hash.Generate(reset.Password, c, e)

	redisClient := RedisClient.NewRedisClient()
	defer redisClient.Close()

	key := User.Username + "_password_token"
	storedToken, err := redisClient.Get(key).Result()
	if err == redis.Nil || storedToken != reset.Token {
		respondError(w, http.StatusUnauthorized, "Invalid token, may have expired. Please try again")
		return
	} else if err != nil {
		log.Println(err)
		respondGenericError(w)
		return
	}

	select {
	case User.Password = <-c:
		break
	case err = <-e:
		log.Println(err)
		respondGenericError(w)
		return
	}

	customErr = userService.UpdatePassword(&User)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	redisClient.Del(key)

	respondJSON(w, http.StatusOK, nil)
}
