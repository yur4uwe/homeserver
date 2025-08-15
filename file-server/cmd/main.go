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

	"file-server/constants"
	"file-server/pkg/handlers"

	"github.com/joho/godotenv"
)

func main() {
	info, err := os.Stat("../../.env")
	if err != nil {
		log.Fatalf("Error stating .env file: %v", err)
	}

	log.Printf("Using .env file at: %s", info.Name())

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.PageHandler)

	server := &http.Server{
		Addr:    ":" + os.Getenv(constants.AppPortKey),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", os.Getenv(constants.AppPortKey))
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
