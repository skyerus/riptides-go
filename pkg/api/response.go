package api

import (
	"encoding/json"
	"github.com/skyerus/riptides-go/pkg/customError"
	"log"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{})  {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"message": message})
}

func respondGenericError(w http.ResponseWriter)  {
	respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Oops, something went wrong. Please try again later."})
}

func respondBadRequest(w http.ResponseWriter)  {
	respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request"})
}

func respondUnauthorizedRequest(w http.ResponseWriter) {
	respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized request"})
}

func handleError(w http.ResponseWriter, customError customError.Error)  {
	if customError.OriginalError() != nil {
		log.Println(customError.OriginalError())
	}
	respondError(w, customError.Code(), customError.Message())
}
