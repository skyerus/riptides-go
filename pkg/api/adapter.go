package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/skyerus/riptides-go/pkg/user"
	"net/http"
)

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func Auth(service user.Service) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if len(token) < 7 {
				respondBadRequest(w)
				return
			}
			token = token[7:]
			tkn, err := service.VerifyToken(token)

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
}
