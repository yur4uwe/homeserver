package db

import (
	"context"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFolderEntryUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with invalid folder data
	invalidFolder := &folder_schema.UserFolder{}
	_, err := con_copy.CreateFolderEntry(context.Background(), invalidFolder, mockUser.Username)
	assert.Error(t, err, "Creating folder entry with invalid folder data should return an error")
	assert.Contains(t, err.Error(), "invalid folder data", "Error message should indicate invalid folder data")

	// Test with empty owner name
	_, err = con_copy.CreateFolderEntry(context.Background(), mockFolder, "")
	assert.Error(t, err, "Creating folder entry with empty owner name should return an error")
	assert.Contains(t, err.Error(), "owner name cannot be empty", "Error message should indicate invalid folder data")
}

func TestCreateRootFolderEntryUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with invalid folder data
	invalidFolder := &folder_schema.UserFolder{}
	_, err := con_copy.CreateRootFolderEntry(context.Background(), invalidFolder, mockUser.Username)
	assert.Error(t, err, "Creating root folder entry with invalid folder data should return an error")
	assert.Contains(t, err.Error(), "invalid folder data", "Error message should indicate invalid folder data")

	// Test with non-empty path
	mockFolder.Path = "non-empty-path"
	_, err = con_copy.CreateRootFolderEntry(context.Background(), mockFolder, mockUser.Username)
	assert.Error(t, err, "Creating root folder entry with non-empty path should return an error")
	assert.Contains(t, err.Error(), "root folder should not have a path", "Error message should indicate invalid path")
}

func TestQueryFolderByNameUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with empty folder name
	_, err := con_copy.QueryFolderByName(context.Background(), "", mockUser.ID, mockFolder.Path)
	assert.Error(t, err, "Querying folder by empty name should return an error")
	assert.Contains(t, err.Error(), "folder does not exist", "Error message should indicate folder does not exist")

	// Test with invalid owner ID
	_, err = con_copy.QueryFolderByName(context.Background(), mockFolder.Name, 0, mockFolder.Path)
	assert.Error(t, err, "Querying folder with invalid owner ID should return an error")
	assert.Contains(t, err.Error(), "folder does not exist", "Error message should indicate query failure")
}

func TestFolderExistsUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with empty folder name
	exists, err := con_copy.FolderExists(context.Background(), "", mockUser.ID, mockFolder.Path)
	assert.NoError(t, err, "Checking folder existence with empty name should return no error")
	assert.False(t, exists, "Folder should not exist for empty name")

	// Test with invalid owner ID
	exists, err = con_copy.FolderExists(context.Background(), mockFolder.Name, 0, mockFolder.Path)
	assert.NoError(t, err, "Checking folder existence with invalid owner ID should return no error")
	assert.False(t, exists, "Folder should not exist for invalid owner ID")
}

func TestExtractFoldersFromDirectoryUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with empty directory
	_, err := con_copy.ExtractFoldersFromDirectory(context.Background(), nil, mockUser)
	assert.Error(t, err, "Extracting folders from empty directory should return an error")
	assert.Contains(t, err.Error(), "directory cannot be empty", "Error message should indicate query failure")

	_, err = con_copy.ExtractFoldersFromDirectory(context.Background(), []string{}, user_schema.NewUserBuilder())
	assert.Error(t, err, "Extracting folders without id directory should return an error")
	assert.Contains(t, err.Error(), "ID", "Error message should indicate query failure")

	// Test with nil requester
	_, err = con_copy.ExtractFoldersFromDirectory(context.Background(), []string{"testuser"}, nil)
	assert.Error(t, err, "Extracting folders with nil requester should return an error")
	assert.Contains(t, err.Error(), "requester cannot be nil", "Error message should indicate invalid requester")
}

func TestDeleteFolderUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with invalid folder data
	invalidFolder := &folder_schema.UserFolder{}
	id, err := con_copy.DeleteFolder(context.Background(), invalidFolder, mockUser.Username)
	assert.Equal(t, 0, id, "Deleting folder with invalid folder data should return ID 0")
	assert.Error(t, err, "Deleting folder with invalid folder data should return an error")
	assert.Contains(t, err.Error(), "invalid folder data", "Error message should indicate invalid folder data")

	// Test with empty owner name
	id, err = con_copy.DeleteFolder(context.Background(), mockFolder, "")
	assert.Equal(t, 0, id, "Deleting folder with empty owner name should return ID 0")
	assert.Error(t, err, "Deleting folder with empty owner name should return an error")
	assert.Contains(t, err.Error(), "invalid folder data", "Error message should indicate invalid folder data")
}

func TestRenameFolderUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with invalid folder data
	invalidFolder := &folder_schema.UserFolder{}
	err := con_copy.RenameFolder(context.Background(), invalidFolder, "newName", mockUser.Username)
	assert.Error(t, err, "Renaming folder with invalid folder data should return an error")
	assert.Contains(t, err.Error(), "invalid folder data", "Error message should indicate invalid folder data")

	// Test with empty new name
	err = con_copy.RenameFolder(context.Background(), mockFolder, "", mockUser.Username)
	assert.Error(t, err, "Renaming folder with empty new name should return an error")
	assert.Contains(t, err.Error(), "invalid folder data", "Error message should indicate invalid folder data")

	// Test with empty owner name
	err = con_copy.RenameFolder(context.Background(), mockFolder, "newName", "")
	assert.Error(t, err, "Renaming folder with empty owner name should return an error")
	assert.Contains(t, err.Error(), "invalid folder data", "Error message should indicate invalid folder data")
}

func TestExtractUserFoldersUnhappyPath(t *testing.T) {
	con_copy, cleanup := setupFolderTestEnvironment(t)
	defer cleanupFolderTestEnvironment(t, cleanup)

	// Test with invalid user ID
	_, err := con_copy.ExtractUserFolders(context.Background(), 0)
	assert.Error(t, err, "Extracting user folders with invalid user ID should return an error")
	assert.Contains(t, err.Error(), "invalid user ID", "Error message should indicate invalid user ID")

	// Test with non-existent user ID
	entries, err := con_copy.ExtractUserFolders(context.Background(), 5)
	assert.NoError(t, err, "Extracting user folders with non-existent user ID should return an error")
	assert.Empty(t, entries, "Extracted folders should be empty for non-existent user ID")
}
