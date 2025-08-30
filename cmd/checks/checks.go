package main

import (
	"homeserver/constants"
	"homeserver/db"
	"homeserver/db/repositories"
	"homeserver/pkg/checks"
	"homeserver/pkg/utils"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	constants.Rootd = utils.GetProjectRoot()

	if err := godotenv.Load(constants.Rootd+".env", constants.Rootd+".env.test"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := checks.CheckVariables(); err != nil {
		log.Fatalf("Environment setup error: %v", err)
	}

	if err := checks.CheckStoragePath(); err != nil {
		log.Fatalf("Storage path setup error: %v", err)
	}

	if err := checks.CheckTestVariables(); err != nil {
		log.Fatalf("Test environment setup error: %v", err)
	}

	if err := checks.CheckTestStoragePath(); err != nil {
		log.Fatalf("Test storage path setup error: %v", err)
	}

	connection := db.NewPool()
	defer connection.Close()

	repos := repositories.NewRepositories(repositories.InitRepositories(connection))

	if err := checks.CheckDatabaseTables(repos); err != nil {
		log.Fatalf("Database table setup error: %v", err)
	}

	if err := checks.CheckDatabaseExtensions(connection); err != nil {
		log.Fatalf("Database extension setup error: %v", err)
	}

	if err := checks.CheckUserFolders(repos); err != nil {
		log.Fatalf("User folder setup error: %v", err)
	}
}
