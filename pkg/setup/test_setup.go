package setup

import (
	"fmt"
	"homeserver/constants"
	"os"
)

const (
	TestDBUrl       = "TEST_DB_URL"
	TestStoragePath = "TEST_STORAGE_PATH"
)

var requiredTestVars = []string{
	TestDBUrl,
}

func CheckTestVariables() error {
	missingVars := make([]string, 0, len(requiredTestVars))

	for _, v := range requiredTestVars {
		if value := os.Getenv(v); value == "" {
			missingVars = append(missingVars, v)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required test environment variables: %v", missingVars)
	}

	fmt.Println("All required test environment variables are set.")

	return nil
}

func CheckTestStoragePath() error {
	value := os.Getenv(constants.StoragePath)

	if value == "" {
		return fmt.Errorf("missing required test environment variable: %s", constants.StoragePath)
	}

	if info, err := os.Stat(value); os.IsNotExist(err) {
		return fmt.Errorf("test storage path does not exist: %s", value)
	} else if !info.IsDir() {
		return fmt.Errorf("test storage path is not a directory: %s", value)
	}

	if value[len(value)-1] != '/' {
		return fmt.Errorf("test storage path must end with a slash: %s", value)
	}

	fmt.Println("Test storage path is set.")

	return nil
}
