package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Starting server on :8080")

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Error starting server: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10* time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Error shutting down server: ", err)
	} else{
		log.Println("Server gracefully stopped")
	}

	log.Println("Server stopped")
}