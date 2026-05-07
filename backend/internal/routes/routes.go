package routes

import "net/http"

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	//Health check route
	mux.HandleFunc("GET /api/health-check-http", handleHealthCheckHTTP)

	return mux
}