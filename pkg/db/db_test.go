package db

import (
	"context"
	"testing"

	"homeserver/pkg/db/file"
	"homeserver/pkg/db/folder"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"homeserver/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestClearEntries(t *testing.T) {
	con_copy, end_transaction := setupTableTestEnvironment(t)
	defer cleanupTableTestEnvironment(t, end_transaction)

	mockOS := new(MockOS)
	mockOS.On("RemoveAll", mock.Anything).Return(nil)
	mockOS.On("Remove", mock.Anything).Return(nil)

	cleanupRemove := test.MockFunction(&file.Remove, mockOS.Remove, func() {})
	defer cleanupRemove()
	cleanupRemoveAll := test.MockFunction(&folder.RemoveAll, mockOS.RemoveAll, func() {})
	defer cleanupRemoveAll()

	err := setup.DatabaseSetup(con_copy)
	assert.NoError(t, err)

	mockUser := user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithEmail("testuser@example.com").
		WithPassword("password123")

	inserted_user, err := con_copy.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err)
	assert.NotNil(t, inserted_user, "Inserted user should not be nil")
	assert.NotEqual(t, 0, inserted_user.ID, "Inserted user ID should not be zero")

	mockFolder := folder_schema.NewUserFolderBuilder().
		WithName(inserted_user.Username).
		WithOwnerID(inserted_user.ID)

	inserted_folder, err := con_copy.CreateRootFolderEntry(context.Background(), mockFolder, inserted_user.Username)
	assert.NoError(t, err, "Creating root folder entry should not return an error")
	assert.NotNil(t, inserted_folder, "Inserted folder should not be nil")
	assert.NotEqual(t, 0, inserted_folder.ID, "Inserted folder ID should not be zero")

	mockFile := file_schema.NewUserFileBuilder().
		WithFilename("testfile").
		WithExtension("txt").
		WithSize(1024).
		WithContentType("text/plain").
		WithPath("testuser").
		WithOwnerID(1)

	_, err = con_copy.CreateFileEntry(context.Background(), mockFile, inserted_user)
	assert.NoError(t, err)

	var countBefore int
	err = con_copy.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM user_files").Scan(&countBefore)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, countBefore)

	_, err = con_copy.ClearEntries(context.Background())
	assert.NoError(t, err)

	var countAfter int
	err = con_copy.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM user_files").Scan(&countAfter)
	assert.NoError(t, err)
	assert.Equal(t, 0, countAfter)
}
