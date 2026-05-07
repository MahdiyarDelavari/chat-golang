package main

import (
	"backend/internal/config"
	"backend/internal/db"
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

	db.InitDB(cfg.DBPath,cfg.DBName)
	defer db.CloseDB()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health-check-http", handleHealthCheckHTTP)

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: mux,
	}

	shutdownCh:= make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt,syscall.SIGTERM,syscall.SIGINT)

	go func () {

		log.Printf("server is running at http://%s",cfg.HTTPServer.Address)

		log.Printf("Health Check Endpoint: http://%s/api/health-check-http",cfg.HTTPServer.Address)

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

func handleHealthCheckHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is healthy:)"))
}