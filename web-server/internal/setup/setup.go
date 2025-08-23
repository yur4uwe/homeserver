package setup

import (
	"homeserver/constants"
	"homeserver/db"
	"homeserver/db/repositories"
	"homeserver/web-server/pkg/setup"
	"homeserver/web-server/pkg/utils"
	"log"

	"github.com/joho/godotenv"
)

func Setup() (*repositories.Repositories, func(), error) {
	constants.Rootd = utils.GetProjectRoot()

	if err := utils.LoadAvailableIconTypes(); err != nil {
		return nil, nil, err
	}

	if err := godotenv.Load(constants.Rootd + ".env"); err != nil {
		return nil, nil, err
	}

	connection, err := db.NewConnection()
	if err != nil {
		return nil, nil, err
	}

	repos := repositories.NewRepositories(repositories.InitRepositories(connection))

	if err := setup.CheckDatabaseExtensions(connection); err != nil {
		log.Fatalf("Database extension check error: %v", err)
		return nil, nil, err
	}

	return repos, connection.Close, nil
}
