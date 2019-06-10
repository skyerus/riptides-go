package api

import (
	"github.com/dgrijalva/jwt-go"
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
		token := r.Header.Get("Authorization")
		if len(token) < 7 {
			respondBadRequest(w)
			return
		}
		token = token[7:]
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
