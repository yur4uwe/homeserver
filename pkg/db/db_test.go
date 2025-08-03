package db

import (
	"context"
	"os"
	"testing"

	"homeserver/pkg/globals"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseFunctionsWithRollback(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	assert.NoError(t, err)

	// Set up a test database connection
	dbURL := os.Getenv("TEST_DB_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	assert.NoError(t, err)
	defer pool.Close()

	database := &Database{pool: pool}

	t.Run("TableExists returns true for existing table", func(t *testing.T) {
		// Begin a transaction
		tx, err := pool.Begin(context.Background())
		assert.NoError(t, err)
		defer tx.Rollback(context.Background()) // Rollback changes after the test

		// Create a test table within the transaction
		_, err = tx.Exec(context.Background(), "CREATE TABLE test_table (id SERIAL PRIMARY KEY)")
		assert.NoError(t, err)

		// Test TableExists
		exists, err := database.TableExists(context.Background(), "test_table")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CreateTable creates a new table", func(t *testing.T) {
		// Begin a transaction
		tx, err := pool.Begin(context.Background())
		assert.NoError(t, err)
		defer tx.Rollback(context.Background()) // Rollback changes after the test

		// Define a mock table
		mockTable := &globals.DBTable{
			Name:        "new_test_table",
			SQLCreation: "CREATE TABLE new_test_table (id SERIAL PRIMARY KEY)",
		}

		// Test CreateTable
		err = database.CreateTable(context.Background(), mockTable)
		assert.NoError(t, err)

		// Verify table creation
		var exists bool
		err = tx.QueryRow(context.Background(), `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)`, "new_test_table").Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}
