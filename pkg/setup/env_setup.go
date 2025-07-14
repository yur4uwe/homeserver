package setup

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"os"
	"strings"
)

const (
	EnvPort          = "PORT"
	EnvDbURL         = "DB_URL"
	EnvStoragePath   = "STORAGE_PATH"
	EnvEmailPassword = "EMAIL_PASSWORD"
	EnvEmailHost     = "EMAIL_HOST"
	EnvEmailPort     = "EMAIL_PORT"
	EnvEmailUser     = "EMAIL_USER"
)

const (
	DBUsersTable = "users"
)

var requiredVars = []string{
	EnvPort,
	EnvDbURL,
	EnvStoragePath,
	EnvEmailPassword,
	EnvEmailHost,
	EnvEmailPort,
	EnvEmailUser,
}

var required_db_tables = []string{
	DBUsersTable,
}

func CheckVariables() error {
	missingVars := make([]string, 0, len(requiredVars))

	for _, v := range requiredVars {
		if value := os.Getenv(v); value == "" {
			missingVars = append(missingVars, v)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	fmt.Println("All required environment variables are set.")

	return nil
}

func CheckStoragePath() error {
	storagePath := os.Getenv(EnvStoragePath)
	if storagePath == "" {
		return fmt.Errorf("storage path is not set in environment variables")
	}

	info, err := os.Stat(storagePath)

	if os.IsNotExist(err) {
		return fmt.Errorf("storage path does not exist: %s", storagePath)
	}

	if !info.IsDir() {
		return fmt.Errorf("storage path is not a directory: %s", storagePath)
	}

	fmt.Println("Storage path exists:", storagePath)
	return nil
}

func CheckDatabaseTables(connection *db.Database) error {
	ctx := context.Background()

	missing_tables := make([]string, 0, len(required_db_tables))
	for _, table := range required_db_tables {
		exists, err := connection.TableExists(ctx, table)
		if err != nil || !exists {
			missing_tables = append(missing_tables, table)
		}
	}

	if len(missing_tables) > 0 {
		return fmt.Errorf("missing required database tables: %s", strings.Join(missing_tables, ", "))
	}

	fmt.Println("All required database tables exist.")
	return nil
}
