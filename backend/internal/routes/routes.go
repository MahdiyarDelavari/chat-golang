package routes

import (
	"backend/internal/middlewares"
	"backend/internal/realtime"
	"net/http"
)

func RegisterRoutes(hub *realtime.Hub) *http.ServeMux {
	mux := http.NewServeMux()

	//Health check route
	mux.HandleFunc("GET /api/health-check-http", handleHealthCheckHTTP)
	mux.HandleFunc("GET /api/health-check-ws", handleHealthCheckWs)

	//Auths
	mux.HandleFunc("POST /api/auth/register-email",handlerEmailRegister) 
	mux.HandleFunc("POST /api/auth/login-email",handlerEmailLogin)
	mux.Handle("POST /api/auth/logout",middlewares.Authenticate(http.HandlerFunc(handlerLogout)))
	mux.HandleFunc("POST /api/auth/refresh-session",handlerRefreshSession)
	mux.Handle("GET /api/auth/current-user",middlewares.Authenticate(http.HandlerFunc(handlerGetCurrentUser)))

	//Users
	mux.Handle("GET /api/users/{id}",middlewares.Authenticate(http.HandlerFunc(handlerGetUserByID)))

	//Conversations
	mux.Handle("GET /api/conversations/privates/{private_id}",middlewares.Authenticate(http.HandlerFunc(handlerGetPrivate)))
	mux.Handle("POST /api/conversations/privates/join",middlewares.Authenticate(http.HandlerFunc(handlerJoinPrivate)))
	mux.Handle("GET /api/conversations",middlewares.Authenticate(http.HandlerFunc(handlerGetConversations)))
	mux.Handle("GET /api/conversations/privates/{private_id}/messages",middlewares.Authenticate(http.HandlerFunc(handlerGetPrivateMessages)))


	//Files
	mux.Handle("POST /api/files/{private_id}",middlewares.Authenticate(http.HandlerFunc(handlerFileUpload)))
	mux.Handle("GET /api/files/", middlewares.Authenticate(handlerGetFile()))

	//Websocket
	mux.HandleFunc("GET /api/ws", func(w http.ResponseWriter, r *http.Request) {handleWebsocket(hub, w, r)})

	return mux
}