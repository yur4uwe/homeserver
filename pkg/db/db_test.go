package db

import (
	"context"
	"testing"

	"homeserver/pkg/globals"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseFunctions(t *testing.T) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	t.Run("TableExists returns true for existing table", func(t *testing.T) {
		// Begin a transaction
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		// Create a test table within the transaction
		_, err := connection.pool.Exec(context.Background(), "CREATE TABLE test_table (id SERIAL PRIMARY KEY)")
		assert.NoError(t, err)

		// Test TableExists
		exists, err := connection.TableExists(context.Background(), "test_table")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CreateTable creates a new table", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		mockTable := &globals.DBTable{
			Name:        "new_test_table",
			SQLCreation: "CREATE TABLE new_test_table (id SERIAL PRIMARY KEY)",
		}

		err := connection.CreateTable(context.Background(), mockTable)
		assert.NoError(t, err)

		var exists bool
		err = connection.pool.QueryRow(context.Background(), `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)`, "new_test_table").Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}
