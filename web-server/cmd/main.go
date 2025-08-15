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

	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/handlers"
	"homeserver/pkg/handlers/auth"
	"homeserver/pkg/middleware"
	"homeserver/pkg/setup"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connection, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer connection.Close()

	globals.FillPages()

	if err := setup.CheckDatabaseExtensions(connection); err != nil {
		log.Fatalf("Database extension check error: %v", err)
	}

	fs := http.FileServer(http.Dir("../static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/verification", auth.VerificationHandler(connection))
	mux.HandleFunc("/upload", handlers.UploadPageHandler)
	mux.HandleFunc("/settings", handlers.SettingsPageHandler)
	mux.HandleFunc("/search/", handlers.SearchRouter(connection))
	mux.HandleFunc("/api/v1/upload", handlers.UploadRouter(connection))
	mux.HandleFunc("/api/v1/login", auth.LoginHandler(connection))
	mux.HandleFunc("/api/v1/register", auth.RegisterHandler(connection))
	mux.HandleFunc("/api/v1/logout", auth.LogoutHandler(connection))
	mux.HandleFunc("/api/v1/file-actions/", handlers.FileActionsRouter(connection))
	mux.HandleFunc("/api/v1/account/settings", handlers.AccountSettingsRouter(connection))
	mux.HandleFunc("/", handlers.StaticPageHandler)

	handler := middleware.MiddlewareWrapper(mux, connection)

	server := &http.Server{
		Handler: handler,
		Addr:    ":" + os.Getenv("PORT"),
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", os.Getenv("MAIN_PORT"))
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
