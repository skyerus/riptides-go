package GoogleHandler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/handler"
	"github.com/skyerus/riptides-go/pkg/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type googleHandler struct {
	redisClient *redis.Client
}

func NewGoogleHandler(redisClient *redis.Client) handler.Handler {
	return &googleHandler{redisClient}
}

func (g googleHandler) SaveCredentials(response *http.Response, user *models.User) customError.Error {
	var gat models.GoogleAccessToken
	err := json.NewDecoder(response.Body).Decode(&gat)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	g.redisClient.Set("gat", gat.AccessToken, time.Duration(gat.ExpiresIn) * time.Second)

	return nil
}

func (g googleHandler) HandleAuthorizedRequest(r *http.Request, user *models.User) customError.Error {
	gat, err := g.redisClient.Get("gat").Result()
	if err != nil {
		if err == redis.Nil {
			return customError.NewHttpError(-1, "", nil)
		}
		return customError.NewGenericHttpError(err)
	}

	r.Header.Set("Authorization", "Bearer " + gat)

	return nil
}

func (g googleHandler) GetRefreshRequest(user *models.User) (*http.Request, customError.Error) {
	var tokenString string
	var googleCreds models.GoogleCredentials
	var request *http.Request

	googleJsonFile, err := os.Open(os.Getenv("GOOGLE_JSON_PATH") + "/google.json")
	if err != nil {
		return request, customError.NewGenericHttpError(err)
	}

	googleJson, err := ioutil.ReadAll(googleJsonFile)
	if err != nil {
		return request, customError.NewGenericHttpError(err)
	}

	err = json.Unmarshal(googleJson, &googleCreds)
	if err != nil {
		return request, customError.NewGenericHttpError(err)
	}

	now := time.Now().Unix()
	claims := models.GoogleClaims{
		Scope: "https://www.googleapis.com/auth/devstorage.full_control",
		StandardClaims: jwt.StandardClaims{
			Issuer: googleCreds.ClientEmail,
			Audience: "https://www.googleapis.com/oauth2/v4/token",
			ExpiresAt: now + 120,
			IssuedAt: now,
		},
	}
	jwtKey := []byte(googleCreds.PrivateKey)
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(jwtKey)
	if err != nil {
		return request, customError.NewGenericHttpError(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err = token.SignedString(parsedKey)
	if err != nil {
		return request, customError.NewGenericHttpError(err)
	}

	body := url.Values{}
	body.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	body.Add("assertion", tokenString)

	request, err = http.NewRequest("POST", "https://www.googleapis.com/oauth2/v4/token", strings.NewReader(body.Encode()))
	if err != nil {
		return request, customError.NewGenericHttpError(err)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}
