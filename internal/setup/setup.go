package setup

import (
	"homeserver/constants"
	"homeserver/db"
	"homeserver/db/repositories"
	"homeserver/pkg/checks"
	"homeserver/pkg/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Setup(loadIcons bool) (*repositories.Repositories, func(), error) {
	constants.Rootd = utils.GetProjectRoot()

	if loadIcons {
		if err := utils.LoadAvailableIconTypes(); err != nil {
			return nil, nil, err
		}
	}

	info, err := os.Stat(constants.Rootd + ".env")
	if err != nil {
		log.Fatalf("Error stating .env file: %v", err)
	}

	log.Printf("Using .env file at: %s", info.Name())

	if err := godotenv.Load(constants.Rootd + ".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connection := db.NewPool()

	repos := repositories.NewRepositories(repositories.InitRepositories(connection))

	if err := checks.CheckDatabaseExtensions(connection); err != nil {
		log.Fatalf("Database extension check error: %v", err)
		return nil, nil, err
	}

	return repos, connection.Close, nil
}
