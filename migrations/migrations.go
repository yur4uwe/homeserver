package main

import (
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

	if err := CreateFolderStructure(connection); err != nil {
		log.Fatalf("Error creating folder structure: %v", err)
	}
}
