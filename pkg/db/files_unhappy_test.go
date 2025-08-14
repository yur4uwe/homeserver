package db

import (
	"context"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFileEntryUnhappyPath(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	invalidFile := &file_schema.UserFile{}
	_, err := con_copy.CreateFileEntry(context.Background(), invalidFile, mockRequester)
	assert.Error(t, err, "Creating file entry with invalid file data should return an error")
	assert.Contains(t, err.Error(), "invalid file data", "Error message should indicate invalid file data")

	invalidRequester := &user_schema.User{}
	_, err = con_copy.CreateFileEntry(context.Background(), mockFile, invalidRequester)
	assert.Error(t, err, "Creating file entry with nil requester should return an error")
	assert.Contains(t, err.Error(), "folder does not exist", "Error message should indicate invalid requester")
}

func TestExtractUserFilesUnhappyPath(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	// Test with invalid user ID
	_, err := con_copy.ExtractUserFiles(context.Background(), 0)
	assert.Error(t, err, "Extracting user files with invalid user ID should return an error")
	assert.Contains(t, err.Error(), "invalid user ID", "Error message should indicate invalid user ID")
}

func TestRenameFileUnhappyPath(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	oldFile := file_schema.NewUserFileBuilder().
		WithFilename("testfile").
		WithExtension("txt").
		WithSize(1024).
		WithContentType("text/plain").
		WithPath(mockRequester.Username).
		WithOwnerID(mockRequester.ID)

	_, err := con_copy.RenameFile(context.Background(), &file_schema.UserFile{}, "newfile.txt", mockRequester)
	assert.Error(t, err, "Renaming file with empty oldPath should return an error")
	assert.Contains(t, err.Error(), "invalid old file data:", "Error message should indicate invalid file path")

	_, err = con_copy.RenameFile(context.Background(), oldFile, "", mockRequester)
	assert.Error(t, err, "Renaming file with empty newPath should return an error")
	assert.Contains(t, err.Error(), "invalid file name", "Error message should indicate invalid file name")

	invalid_requester := &user_schema.User{}
	_, err = con_copy.RenameFile(context.Background(), oldFile, "newfile.txt", invalid_requester)
	assert.Error(t, err, "Renaming file with invalid requester should return an error")
	assert.Contains(t, err.Error(), "error validating requester", "Error message should indicate invalid requester")
}

func TestDeleteFileUnhappyPath(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	// Test with invalid file data
	invalidFile := &file_schema.UserFile{}
	_, err := con_copy.DeleteFile(context.Background(), invalidFile, mockRequester.Username)
	assert.Error(t, err, "Deleting file with invalid file data should return an error")
	assert.Contains(t, err.Error(), "invalid file data", "Error message should indicate invalid file data")

	// Test with empty username
	_, err = con_copy.DeleteFile(context.Background(), mockFile, "")
	assert.Error(t, err, "Deleting file with empty username should return an error")
	assert.Contains(t, err.Error(), "user name cannot be empty", "Error message should indicate invalid username")
}

func TestGetFileUnhappyPath(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	// Test with invalid file data
	invalidFile := &file_schema.UserFile{}
	_, err := con_copy.GetFile(context.Background(), invalidFile, mockRequester)
	assert.Error(t, err, "Getting file with invalid file data should return an error")
	assert.Contains(t, err.Error(), "invalid file data", "Error message should indicate invalid file data")

	invalid_requester := &user_schema.User{}
	_, err = con_copy.GetFile(context.Background(), mockFile, invalid_requester)
	assert.Error(t, err, "Getting file with invalid requester should return an error")
	assert.Contains(t, err.Error(), "error validating requester", "Error message should indicate invalid requester")
}

func TestExtractFilesFromDirectoryUnhappyPath(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	// Test with invalid directory
	files, err := con_copy.ExtractFilesFromDirectory(context.Background(), "asfdas", mockRequester)
	assert.NoError(t, err, "Extracting files from invalid directory should return an error")
	assert.Empty(t, files, "Files should be nil when an error occurs")

	invalid_requester := &user_schema.User{}
	_, err = con_copy.ExtractFilesFromDirectory(context.Background(), "testuser", invalid_requester)
	assert.Error(t, err, "Extracting files with invalid requester should return an error")
	assert.Contains(t, err.Error(), "requester must have a valid user ID", "Error message should indicate invalid requester")
}

func TestSaveFileUnhappyPath(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	// Test with invalid file header
	invalidHeader := &file_schema.UserFile{}
	_, _, err := con_copy.SaveFile(context.Background(), invalidHeader, nil, mockRequester)
	assert.Error(t, err, "Saving file with invalid file header should return an error")
	assert.Contains(t, err.Error(), "file header must have valid fields", "Error message should indicate invalid file header")

	// Test with nil file data
	assert.Panics(t, func() {
		_, _, err = con_copy.SaveFile(context.Background(), mockFile, nil, mockRequester)
	}, "Saving file with nil file data should panic")

	invalid_requester := &user_schema.User{}
	_, _, err = con_copy.SaveFile(context.Background(), mockFile, nil, invalid_requester)
	assert.Error(t, err, "Saving file with invalid requester should return an error")
	assert.Contains(t, err.Error(), "requester must have a valid user ID", "Error message should indicate invalid requester")
}
