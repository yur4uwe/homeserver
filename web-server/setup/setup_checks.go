package main

import (
	"homeserver/pkg/db"
	"homeserver/pkg/setup"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env", "../.env.test"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := setup.CheckVariables(); err != nil {
		log.Fatalf("Environment setup error: %v", err)
	}

	if err := setup.CheckStoragePath(); err != nil {
		log.Fatalf("Storage path setup error: %v", err)
	}

	if err := setup.CheckTestVariables(); err != nil {
		log.Fatalf("Test environment setup error: %v", err)
	}

	if err := setup.CheckTestStoragePath(); err != nil {
		log.Fatalf("Test storage path setup error: %v", err)
	}

	connection, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer connection.Close()

	if err := setup.CheckDatabaseTables(connection); err != nil {
		log.Fatalf("Database table setup error: %v", err)
	}

	if err := setup.CheckDatabaseExtensions(connection); err != nil {
		log.Fatalf("Database extension setup error: %v", err)
	}

	if err := setup.CheckUserFolders(connection); err != nil {
		log.Fatalf("User folder setup error: %v", err)
	}
}
