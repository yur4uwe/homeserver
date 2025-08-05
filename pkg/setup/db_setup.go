package setup

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"os"
)

var RequiredDBTables = []*globals.DBTable{
	globals.DBUsersTable,
	globals.DBFoldersTable,
	globals.DBFilesTable,
}

func CheckDatabaseTables(connection interfaces.DatabaseInterface, enable_logs ...bool) error {
	ctx := context.Background()

	to_log := len(enable_logs) > 0 && enable_logs[0]

	for _, table := range RequiredDBTables {
		exists, err := connection.TableExists(ctx, table.Name)
		if err != nil {
			return fmt.Errorf("error checking table %s: %v", table.Name, err)
		}

		if to_log {
			fmt.Printf("Checking table %s: exists=%v\n", table.Name, exists)
		}

		if !exists {
			if err := connection.CreateTable(ctx, table); err != nil {
				return fmt.Errorf("error creating table %s: %v", table.Name, err)
			}
		}
	}

	if to_log {
		fmt.Println("All required database tables exist.")
	}
	return nil
}

func CheckUserFolders(connection interfaces.DatabaseInterface, enable_logs ...bool) error {
	ctx := context.Background()

	to_log := len(enable_logs) > 0 && enable_logs[0]

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

	if to_log {
		fmt.Println("All user folders have been created.")
	}
	return nil
}

func CheckDatabaseExtensions(connection interfaces.DatabaseInterface) error {
	if _, err := connection.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS pg_trgm;"); err != nil {
		return fmt.Errorf("error creating extension: %v", err)
	}

	return nil
}

func CheckTableIndices(connection interfaces.DatabaseInterface) error {
	if _, err := connection.Exec(context.Background(),
		`CREATE INDEX IF NOT EXISTS user_files_filename_trgm_idx 
		ON user_files 
		USING gin (filename gin_trgm_ops);`,
	); err != nil {
		return fmt.Errorf("error creating index on user_files table: %v", err)
	}

	return nil
}

func DatabaseSetup(connection interfaces.DatabaseInterface, enable_logs ...bool) error {
	if err := CheckDatabaseExtensions(connection); err != nil {
		return fmt.Errorf("database extension setup failed: %v", err)
	}

	if err := CheckDatabaseTables(connection, enable_logs...); err != nil {
		return fmt.Errorf("database setup failed: %v", err)
	}

	if err := CheckTableIndices(connection); err != nil {
		return fmt.Errorf("table index setup failed: %v", err)
	}

	if err := CheckUserFolders(connection, enable_logs...); err != nil {
		return fmt.Errorf("user folder setup failed: %v", err)
	}

	if len(enable_logs) > 0 && enable_logs[0] {
		fmt.Println("Database and user folders are set up successfully.")
	}
	return nil
}
