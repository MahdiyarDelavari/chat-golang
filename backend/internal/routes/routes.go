package routes

import (
	"backend/internal/middlewares"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	//Health check route
	mux.HandleFunc("GET /api/health-check-http", handleHealthCheckHTTP)

	//Auths
	mux.HandleFunc("POST /api/auth/register-email",handlerEmailRegister) 
	mux.HandleFunc("POST /api/auth/login-email",handlerEmailLogin)
	mux.Handle("POST /api/auth/logout",middlewares.Authenticate(http.HandlerFunc(handlerLogout)))
	mux.HandleFunc("POST /api/auth/refresh-session",handlerRefreshSession)
	mux.Handle("GET /api/auth/current-user",middlewares.Authenticate(http.HandlerFunc(handlerGetCurrentUser)))

	return mux
}