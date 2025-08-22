package main

import (
	"log"

	"homeserver/web-server/internal/routes"
	"homeserver/web-server/internal/server"
	"homeserver/web-server/internal/setup"
)

func main() {
	repos, cleanup, err := setup.Setup()
	if err != nil {
		log.Fatalf("Setup error: %v", err)
	}
	defer cleanup()

	handler := routes.NewRouter(repos)

	if err := server.Run(handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
