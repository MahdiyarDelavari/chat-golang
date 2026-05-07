package routes

import "net/http"

func handleHealthCheckHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is healthy:)"))
}