package checks

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/constants/dbconsts"
	"homeserver/db/repositories"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var RequiredDBTables = []*dbconsts.DBTable{
	dbconsts.DBUsersTable,
	dbconsts.DBFoldersTable,
	dbconsts.DBFilesTable,
}

func CheckDatabaseTables(repos *repositories.Repositories, enable_logs ...bool) error {
	ctx := context.Background()

	to_log := len(enable_logs) > 0 && enable_logs[0]

	for _, table := range RequiredDBTables {
		exists, err := repos.TableRepo.TableExists(ctx, table.Name)
		if err != nil {
			return fmt.Errorf("error checking table %s: %v", table.Name, err)
		}

		if to_log {
			fmt.Printf("Checking table %s: exists=%v\n", table.Name, exists)
		}

		if !exists {
			if err := repos.TableRepo.CreateTable(ctx, table); err != nil {
				return fmt.Errorf("error creating table %s: %v", table.Name, err)
			}
		}
	}

	if to_log {
		fmt.Println("All required database tables exist.")
	}
	return nil
}

func CheckUserFolders(repos *repositories.Repositories, enable_logs ...bool) error {
	ctx := context.Background()

	to_log := len(enable_logs) > 0 && enable_logs[0]

	storage_path := os.Getenv(constants.StoragePath)

	users, err := repos.UserRepo.QueryAll(ctx)
	if err != nil {
		return fmt.Errorf("error querying users: %v", err)
	}

	for _, user := range users {
		if err := os.MkdirAll(fmt.Sprintf("%s%s", storage_path, user.Username), os.ModePerm); err != nil {
			return fmt.Errorf("error creating folder for user %s: %v", user.Username, err)
		}

		rootFolder := folschema.New().
			WithName(user.Username).
			WithOwnerID(user.ID)

		if folder, err := repos.FolderRepo.GetFolder(ctx, rootFolder, &user); err != nil {
			return fmt.Errorf("error getting folder for user %s: %v", user.Username, err)
		} else if folder == nil {
			if _, err := repos.FolderRepo.CreateFolder(ctx, rootFolder, &user); err != nil {
				return fmt.Errorf("error creating folder for user %s: %v", user.Username, err)
			}
		}
	}

	if to_log {
		fmt.Println("All user folders have been created.")
	}
	return nil
}

func CheckDatabaseExtensions(connection *pgxpool.Pool) error {
	if _, err := connection.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS pg_trgm;"); err != nil {
		return fmt.Errorf("error creating extension: %v", err)
	}

	if _, err := connection.Exec(context.Background(), "SELECT set_limit(0.1);"); err != nil {
		return fmt.Errorf("error setting similarity threshold: %v", err)
	}

	var similarityThreshold string
	err := connection.QueryRow(context.Background(), "SELECT show_limit();").Scan(&similarityThreshold)
	if err != nil {
		return fmt.Errorf("error finding similarity threshold: %v", err)
	}

	fmt.Printf("Current similarity threshold: %s\n", similarityThreshold)

	return nil
}

func CheckTableIndices(connection *pgxpool.Pool) error {
	if _, err := connection.Exec(context.Background(),
		`CREATE INDEX IF NOT EXISTS user_files_filename_trgm_idx 
		ON user_files 
		USING gin (filename gin_trgm_ops);`,
	); err != nil {
		return fmt.Errorf("error creating index on user_files table: %v", err)
	}

	return nil
}

// func DatabaseSetup(repos *repositories.Repositories, enable_logs ...bool) error {
// 	if err := CheckDatabaseExtensions(connection); err != nil {
// 		return fmt.Errorf("database extension setup failed: %v", err)
// 	}

// 	if err := CheckDatabaseTables(connection, enable_logs...); err != nil {
// 		return fmt.Errorf("database setup failed: %v", err)
// 	}

// 	if err := CheckTableIndices(connection); err != nil {
// 		return fmt.Errorf("table index setup failed: %v", err)
// 	}

// 	if err := CheckUserFolders(connection, enable_logs...); err != nil {
// 		return fmt.Errorf("user folder setup failed: %v", err)
// 	}

// 	if len(enable_logs) > 0 && enable_logs[0] {
// 		fmt.Println("Database and user folders are set up successfully.")
// 	}
// 	return nil
// }
