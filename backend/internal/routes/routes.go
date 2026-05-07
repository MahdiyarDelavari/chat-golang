package routes

import "net/http"

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	//Health check route
	mux.HandleFunc("GET /api/health-check-http", handleHealthCheckHTTP)

	//Auths
	mux.HandleFunc("POST /api/auth/register-email",handlerEmailRegister) 

	return mux
}