package routes

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {

	//Health check route
	mux.HandleFunc("GET /api/health-check-http", handleHealthCheckHTTP)
	

}