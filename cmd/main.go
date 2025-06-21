package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homeserver/pkg/handlers"
	"homeserver/pkg/middleware"
	"homeserver/pkg/setup"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := setup.CheckVariables(); err != nil {
		log.Fatalf("Environment setup error: %v", err)
	}

	if err := setup.CheckStoragePath(); err != nil {
		log.Fatalf("Storage path setup error: %v", err)
	}

	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Hello)

	handler := middleware.MiddlewareWrapper(mux)

	server := &http.Server{
		Handler: handler,
		Addr:    ":" + port,
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}
