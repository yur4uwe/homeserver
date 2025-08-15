package main

import (
	"homeserver/migrations/scripts"
	"homeserver/pkg/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connection, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer connection.Close()

	// if err := database.CreateFolderStructure(connection); err != nil {
	// 	log.Fatalf("Error creating folder structure: %v", err)
	// }

	// if err := database.SeparateExtensionsFromNames(connection); err != nil {
	// 	log.Fatalf("Error separating extensions from names: %v", err)
	// }

	// if err := database.ResetFolderStructure(connection); err != nil {
	// 	log.Fatalf("Error resetting folder structure: %v", err)
	// }

	scripts.ChangeClassNames()
}
