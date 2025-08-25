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

	"homeserver/constants"
	"homeserver/db"
	"homeserver/db/repositories"
	"homeserver/file-server/internal/logger"
	"homeserver/file-server/pkg/handlers"
	"homeserver/file-server/pkg/middleware"
	"homeserver/web-server/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {
	constants.Rootd = utils.GetProjectRoot()

	info, err := os.Stat(constants.Rootd + ".env")
	if err != nil {
		log.Fatalf("Error stating .env file: %v", err)
	}

	log.Printf("Using .env file at: %s", info.Name())

	if err := godotenv.Load(constants.Rootd + ".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connection, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Error establishing database connection: %v", err)
	}

	conn := repositories.NewRepositories(repositories.InitRepositories(connection))

	fs := http.FileServer(http.Dir(constants.Rootd + "static/"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/files/", handlers.FileHandler)
	mux.HandleFunc("/folders/", handlers.FolderHandler)

	handler := middleware.Wrapper(mux, conn, &logger.Logger{})

	server := &http.Server{
		Addr:    ":" + os.Getenv(constants.FileSrvPort),
		Handler: handler,
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", os.Getenv(constants.FileSrvPort))
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
