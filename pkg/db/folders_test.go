package db

import (
	"context"
	"homeserver/pkg/db/folder"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"homeserver/test"
	"homeserver/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupFolderTestEnvironment(t *testing.T) (*Database, func()) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	mockOS := new(mocks.MockOS)
	mockOS.On("MkdirAll", mock.Anything, mock.Anything).Return(nil)
	mockOS.On("RemoveAll", mock.Anything).Return(nil)
	mockOS.On("Rename", mock.Anything, mock.Anything).Return(nil)

	cleanupMkdirAll := test.MockFunction(&folder.MkdirAll, mockOS.MkdirAll, func() {})
	cleanupRemoveAll := test.MockFunction(&folder.RemoveAll, mockOS.RemoveAll, func() {})
	cleanupRename := test.MockFunction(&folder.Rename, mockOS.Rename, func() {})

	err := setup.DatabaseSetup(connection)
	assert.NoError(t, err, "Failed to check database tables")

	mockUser = user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithEmail("testuser@example.com").
		WithPassword("securepassword")

	mockUser, err = connection.CreateUserEntry(context.Background(), mockUser)
	assert.NoError(t, err, "Failed to create mock user")
	assert.NotNil(t, mockUser, "Mock user should not be nil")
	assert.NotEqual(t, 0, mockUser.ID, "Mock user ID should not be zero")

	mockFolder = folder_schema.NewUserFolderBuilder().
		WithName("testFolder").
		WithPath("testuser").
		WithOwnerID(mockUser.ID)

	mockRootFolder = folder_schema.NewUserFolderBuilder().
		WithName("testuser").
		WithOwnerID(mockUser.ID)

	mockRootFolder.Path = ""

	con_copy := connection.Copy()
	end_transaction := con_copy.UseTransaction()
	assert.NotNil(t, end_transaction)

	cleanup := func() {
		cleanupMkdirAll()
		cleanupRemoveAll()
		cleanupRename()
		end_transaction()
	}

	return con_copy, cleanup
}

func cleanupFolderTestEnvironment(t *testing.T, cleanupFunc func()) {
	cleanupFunc()
	for _, table := range setup.RequiredDBTables {
		err := connection.DropTable(context.Background(), table.Name)
		if err != nil {
			t.Errorf("Failed to drop table %s: %v", table.Name, err)
		}
	}
}

func TestCreateFolderEntry(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)
	assert.NotNil(t, con_copy, "Database connection copy should not be nil")

	// Create a folder entry
	inserted_folder, err := con_copy.CreateFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.NoError(t, err, "Creating folder entry should not return an error")

	// Verify the folder exists
	exists, err := con_copy.FolderExists(context.Background(), inserted_folder.Name, mockUser.ID, inserted_folder.Path)
	assert.NoError(t, err, "Checking folder existence should not return an error")
	assert.True(t, exists, "Folder should exist after creation")
}

func TestCreateRootFolderEntry(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Create a root folder entry
	inserted_folder, err := con_copy.CreateRootFolderEntry(context.Background(), mockRootFolder, mockUser.Username)
	assert.NoError(t, err, "Creating root folder entry should not return an error")

	// Verify the folder exists
	exists, err := con_copy.FolderExists(context.Background(), inserted_folder.Name, mockUser.ID, "")
	assert.NoError(t, err, "Checking root folder existence should not return an error")
	assert.True(t, exists, "Root folder should exist after creation")
}

func TestQueryFolderByName(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Create a folder entry
	inserted_folder, err := con_copy.CreateFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.NoError(t, err, "Creating folder entry should not return an error")

	// Query the folder by name
	folder, err := con_copy.QueryFolderByName(context.Background(), inserted_folder.Name, mockUser.ID, inserted_folder.Path)
	assert.NoError(t, err, "Querying folder by name should not return an error")
	assert.NotNil(t, folder, "Folder should not be nil")
	assert.Equal(t, inserted_folder.ID, folder.ID, "Folder ID should match")
}

func TestFolderExists(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Create a folder entry
	inserted_folder, err := con_copy.CreateFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.NoError(t, err, "Creating folder entry should not return an error")

	// Check if the folder exists
	exists, err := con_copy.FolderExists(context.Background(), inserted_folder.Name, mockUser.ID, inserted_folder.Path)
	assert.NoError(t, err, "Checking folder existence should not return an error")
	assert.True(t, exists, "Folder should exist after creation")
}

func TestExtractFoldersFromDirectory(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Create a folder entry
	inserted_folder, err := con_copy.CreateFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.NoError(t, err, "Creating folder entry should not return an error")

	// Extract folders from the directory
	folders, err := con_copy.ExtractFoldersFromDirectory(context.Background(), []string{inserted_folder.Path}, mockUser)
	assert.NoError(t, err, "Extracting folders from directory should not return an error")
	assert.NotEmpty(t, folders, "Folders should not be empty")
	assert.Equal(t, mockFolder.Name, folders[0].Name, "Folder name should match")
}

func TestDeleteFolder(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Create a folder entry
	inserted_folder, err := con_copy.CreateFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.NoError(t, err, "Creating folder entry should not return an error")

	// Delete the folder
	id, err := con_copy.DeleteFolder(context.Background(), inserted_folder, mockUser.Username)
	assert.NoError(t, err, "Deleting folder should not return an error")
	assert.Equal(t, inserted_folder.ID, id, "Deleting folder should return the correct ID")

	// Verify the folder no longer exists
	exists, err := con_copy.FolderExists(context.Background(), inserted_folder.Name, mockUser.ID, inserted_folder.Path)
	assert.NoError(t, err, "Checking folder existence should not return an error")
	assert.False(t, exists, "Folder should not exist after deletion")
}

func TestRenameFolder(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Create a folder entry
	inserted_folder, err := con_copy.CreateFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.NoError(t, err, "Creating folder entry should not return an error")

	// Rename the folder
	newName := "renamedFolder"
	err = con_copy.RenameFolder(context.Background(), inserted_folder, newName, mockUser.Username)
	assert.NoError(t, err, "Renaming folder should not return an error")

	// Verify the folder exists with the new name
	exists, err := con_copy.FolderExists(context.Background(), newName, mockUser.ID, inserted_folder.Path)
	assert.NoError(t, err, "Checking folder existence should not return an error")
	assert.True(t, exists, "Folder should exist with the new name after renaming")
}

func TestExtractUserFolders(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Create a folder entry
	inserted_folder, err := con_copy.CreateFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.NoError(t, err, "Creating folder entry should not return an error")

	// Extract user folders
	folders, err := con_copy.ExtractUserFolders(context.Background(), mockUser.ID)
	assert.NoError(t, err, "Extracting user folders should not return an error")
	assert.NotEmpty(t, folders, "User folders should not be empty")
	assert.Equal(t, inserted_folder.Name, folders[0].Name, "Extracted folder name should match")
}
