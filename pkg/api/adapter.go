package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, customErr := GetAuthToken(r)
		if customErr != nil {
			respondUnauthorizedRequest(w)
			return
		}
		claims := &models.Claims{}

		jwtFile, err := os.Open(os.Getenv("JWT_PATH") + "/private.pem")
		if err != nil {
			log.Fatal(err)
		}
		defer jwtFile.Close()

		jwtKey, err := ioutil.ReadAll(jwtFile)

		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				respondUnauthorizedRequest(w)
				return
			}
			respondUnauthorizedRequest(w)
			return
		}
		if !tkn.Valid {
			respondUnauthorizedRequest(w)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ALLOW_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func GetAuthToken(r *http.Request) (string, customError.Error) {
	token := r.Header.Get("Authorization")
	if len(token) < 7 {
		return token, customError.NewGenericHttpError(nil)
	}
	return token[7:], nil
}
