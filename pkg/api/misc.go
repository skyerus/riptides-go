package api

import "net/http"

func HealthCheck(w http.ResponseWriter, r *http.Request)  {
	respondJSON(w, 200, nil)
}

func ApiHealthCheck(w http.ResponseWriter, r *http.Request)  {
	respondJSON(w, 200, nil)
}
