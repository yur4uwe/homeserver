package server

import (
	"context"
	"errors"
	"homeserver/constants"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(handler http.Handler) error {
	server := &http.Server{
		Handler: handler,
		Addr:    ":" + os.Getenv(constants.WebSrvPort),
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", os.Getenv(constants.WebSrvPort))
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
		return err
	}
	log.Println("Graceful shutdown complete.")

	return nil
}
