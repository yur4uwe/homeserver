package db

import (
	"context"
	"testing"

	"homeserver/pkg/globals"

	"github.com/stretchr/testify/assert"
)

func setupTableTestEnvironment(t *testing.T) (*Database, func()) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)

	return con_copy, EndTransaction
}

func cleanupTableTestEnvironment(t *testing.T, cleanupFunc func()) {
	assert.NotNil(t, connection, "Database connection should not be nil")
	cleanupFunc()
}

func TestTableExists(t *testing.T) {
	con_copy, end_transaction := setupTableTestEnvironment(t)
	defer cleanupTableTestEnvironment(t, end_transaction)

	_, err := con_copy.pool.Exec(context.Background(), "CREATE TABLE test_table (id SERIAL PRIMARY KEY)")
	assert.NoError(t, err)

	exists, err := con_copy.TableExists(context.Background(), "test_table")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCreateTable(t *testing.T) {
	con_copy, end_transaction := setupTableTestEnvironment(t)
	defer cleanupTableTestEnvironment(t, end_transaction)

	mockTable := &globals.DBTable{
		Name:        "new_test_table",
		SQLCreation: "CREATE TABLE new_test_table (id SERIAL PRIMARY KEY)",
	}

	err := con_copy.CreateTable(context.Background(), mockTable)
	assert.NoError(t, err)

	var exists bool
	err = con_copy.pool.QueryRow(context.Background(), `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)`, "new_test_table").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists)
}
