package checks

import (
	"fmt"
	"homeserver/constants"
	"os"
)

var requiredVars = []string{
	constants.WebSrvPort,
	constants.WebSrvHost,
	constants.FileSrvPort,
	constants.FileSrvHost,
	constants.StaticPort,
	constants.StaticHost,
	constants.DbURL,
	constants.StoragePath,
	constants.EmailPassword,
	constants.EmailHost,
	constants.EmailPort,
	constants.EmailUser,
	constants.SourceSecret,
	constants.TokenSecret,
}

func CheckVariables() error {
	varsErrors := map[string]string{}

	for _, v := range requiredVars {
		val := os.Getenv(v)
		if val == "" {
			varsErrors[v] = "missing"
		} else if v == constants.StoragePath {
			if info, err := os.Stat(val); os.IsNotExist(err) || !info.IsDir() {
				varsErrors[v] = "invalid"
			} else if val[len(val)-1] != '/' && val[len(val)-1] != '\\' {
				varsErrors[v] = "must end with a slash"
			}
		}
	}

	if len(varsErrors) > 0 {
		return fmt.Errorf("missing required environment variables: %v", varsErrors)
	}

	fmt.Println("All required environment variables are set.")

	return nil
}

func CheckStoragePath() error {
	storagePath := os.Getenv(constants.StoragePath)
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
