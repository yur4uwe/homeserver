package setup

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"os"
)

var required_db_tables = []*globals.DBTable{
	globals.DBUsersTable,
	globals.DBFilesTable,
	globals.DBFoldersTable,
}

func CheckDatabaseTables(connection *db.Database) error {
	ctx := context.Background()

	for _, table := range required_db_tables {
		exists, err := connection.TableExists(ctx, table.Name)
		if err != nil {
			return fmt.Errorf("error checking table %s: %v", table.Name, err)
		}

		if !exists {
			connection.CreateTable(ctx, table)
		}
	}

	fmt.Println("All required database tables exist.")
	return nil
}

func CheckUserFolders(connection *db.Database) error {
	ctx := context.Background()

	storage_path := os.Getenv(globals.EnvStoragePath)

	users, err := connection.ExtractUserEntries(ctx)
	if err != nil {
		return fmt.Errorf("error extracting user entries: %v", err)
	}

	for _, user := range users {
		err := connection.CheckUserFolder(ctx, user)
		if err != nil {
			return fmt.Errorf("error checking folder for user %s: %v", user.Username, err)
		}

		if info, err := os.Stat(fmt.Sprintf("%s%s", storage_path, user.Username)); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(fmt.Sprintf("%s%s", storage_path, user.Username), os.ModePerm); err != nil {
					return fmt.Errorf("error creating folder for user %s: %v", user.Username, err)
				}
			} else {
				return fmt.Errorf("error checking folder for user %s: %v", user.Username, err)
			}
		} else if !info.IsDir() {
			return fmt.Errorf("storage path for user %s is not a directory", user.Username)
		}
	}

	fmt.Println("All user folders have been created.")
	return nil
}
