package main

import (
	"homeserver/constants"
	"homeserver/db"
	"homeserver/db/repositories"
	"homeserver/pkg/setup"
	"homeserver/pkg/utils"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	constants.Rootd = utils.GetProjectRoot()

	if err := godotenv.Load(constants.Rootd+".env", constants.Rootd+".env.test"); err != nil {
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

	connection, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer connection.Close()

	repos := repositories.NewRepositories(repositories.InitRepositories(connection))

	if err := setup.CheckDatabaseTables(repos); err != nil {
		log.Fatalf("Database table setup error: %v", err)
	}

	if err := setup.CheckDatabaseExtensions(connection); err != nil {
		log.Fatalf("Database extension setup error: %v", err)
	}

	if err := setup.CheckUserFolders(repos); err != nil {
		log.Fatalf("User folder setup error: %v", err)
	}
}
