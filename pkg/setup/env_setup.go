package setup

import (
	"fmt"
	"os"
)

const (
	EnvPort     = "PORT"
	EnvDbURL    = "DB_URL"
	StoragePath = "STORAGE_PATH"
)

var requiredVars = []string{
	EnvPort,
	StoragePath,
	EnvDbURL,
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
	storagePath := os.Getenv(StoragePath)
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
