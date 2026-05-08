package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/middlewares"
	"backend/internal/routes"
	"backend/internal/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	cfg := config.LoadConfig()
	utils.InitJWT(cfg.JWTKey)

	db.InitDB(cfg.DBPath,cfg.DBName)
	defer db.CloseDB()

	mux := routes.RegisterRoutes()

	loggerMux := middlewares.LoggingMiddleware(mux)
	corsMux := middlewares.CorsMiddleware(loggerMux)

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: corsMux,
	}

	shutdownCh:= make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt,syscall.SIGTERM,syscall.SIGINT)

	go func () {
		// Log server start and available endpoints
		log.Printf("server is running at http://%s",server.Addr)

		// Log health check endpoint
		log.Printf("Health Check Endpoint: http://%s/api/health-check-http",server.Addr)
		log.Printf("Health Check WS, GET: ws://%s/api/health-check-ws", server.Addr)

		// Log auth endpoints
		log.Printf("Email register, POST http://%s/api/auth/register-email",server.Addr)
		log.Printf("Email login, POST http://%s/api/auth/login-email",server.Addr)
		log.Printf("Logout, POST http://%s/api/auth/logout (requires auth)",server.Addr)
		log.Printf("Session Refresh, POST http://%s/api/auth/refresh-session (requires auth)",server.Addr)
		log.Printf("Get Current User, GET http://%s/api/auth/current-user (requires auth)",server.Addr)

		//Users
		log.Printf("Get User By Id, GET http://%s/api/users/{user_id} (requires auth)",server.Addr)

		// Conversations
		log.Printf("Get Conversations, GET http://%s/api/conversations/privates/{private_id} (requires auth)",server.Addr)
		log.Printf("Join Conversation, POST http://%s/api/conversations/privates/join (requires auth)",server.Addr)
		log.Printf("Get All Conversations, GET http://%s/api/conversations (requires auth)",server.Addr)
		log.Printf("Get Conversation Messages (Paginated), GET http://%s/api/conversations/privates/{private_id}/messages?page=1&limit=20 (requires auth)",server.Addr)

		// Files
		log.Printf("File Upload, POST http://%s/api/files/{private_id} (requires auth)",server.Addr)
		log.Printf("File Download, GET http://%s/api/files/ (requires auth)",server.Addr)


		err:=server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server: ", err)
		}
	}()

	sig:= <- shutdownCh
	log.Printf("Received signal %s, shutting down server...", sig)


	ctx, cancel := context.WithTimeout(context.Background(), 10* time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Error shutting down server: ", err)
	} else{
		log.Println("Server gracefully stopped")
	}
	
	signal.Stop(shutdownCh)
	close(shutdownCh)

	log.Println("--Server stopped successfully--")
}