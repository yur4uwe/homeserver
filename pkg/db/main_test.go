package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting test suite...")

	err := godotenv.Load("../../.env.test")
	if err != nil {
		panic("Failed to load .env.test file")
	}

	connection, err = NewTestDatabase()
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}
	defer connection.Close()

	code := m.Run()

	os.Exit(code)
}
