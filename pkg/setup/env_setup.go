package setup

import (
	"fmt"
	"os"
)

const (
	EnvPort  = "PORT"
	EnvDbURL = "DB_URL"
)

var requiredVars = []string{
	EnvPort,
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
