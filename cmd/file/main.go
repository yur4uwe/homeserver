package main

import (
	"log"
	"os"

	"homeserver/constants"
	"homeserver/internal/routes"
	"homeserver/internal/server"
	"homeserver/internal/setup"
)

func main() {
	repos, cleanup, err := setup.Setup(false)
	if err != nil {
		log.Fatalf("Setup error: %v", err)
	}
	defer cleanup()

	handler := routes.NewFileRouter(repos)

	if err := server.Run(handler, os.Getenv(constants.FileSrvPort)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
