package db

import (
	"context"
	"homeserver/pkg/db/file"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"homeserver/test"
	"homeserver/test/mocks"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupFileTestEnvironment(t *testing.T) cleanupFunc {
	assert.NotNil(t, connection, "Database connection should not be nil")

	mockOS := new(mocks.MockOS)
	mockOS.On("MkdirAll", mock.Anything, mock.Anything).Return(nil)
	mockOS.On("Rename", mock.Anything, mock.Anything).Return(nil)
	mockOS.On("RemoveAll", mock.Anything).Return(nil)
	mockOS.On("Create", mock.Anything).Return(nil, nil)
	mockOS.On("Remove", mock.Anything).Return(nil)
	mockOS.On("Stat", mock.Anything).Return(nil, os.ErrNotExist)

	cleanupMkdir := test.MockFunction(&file.MkdirAll, mockOS.MkdirAll, func() {})
	cleanupRename := test.MockFunction(&file.Rename, mockOS.Rename, func() {})
	cleanupRemoveAll := test.MockFunction(&file.RemoveAll, mockOS.RemoveAll, func() {})
	cleanupCreate := test.MockFunction(&file.Create, mockOS.Create, func() {})
	cleanupRemove := test.MockFunction(&file.Remove, mockOS.Remove, func() {})
	cleanupStat := test.MockFunction(&file.Stat, mockOS.Stat, func() {})

	mockIO := new(mocks.MockIO)
	mockIO.On("Copy", mock.Anything, mock.Anything).Return(int64(0), nil)

	cleanupCopy := test.MockFunction(&file.Copy, mockIO.Copy, func() {})

	err := setup.DatabaseSetup(connection, false)
	assert.NoError(t, err, "Failed to check database tables")

	mockRequester = user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithEmail("testuser@example.com").
		WithPassword("password123")

	mockRequester, err = connection.CreateUserEntry(context.Background(), mockRequester)
	assert.NoError(t, err, "Failed to create user entry")

	mockFolder = folder_schema.NewUserFolderBuilder().
		WithName(mockRequester.Username).
		WithOwnerID(mockRequester.ID)

	mockRootFolder, err = connection.CreateRootFolderEntry(context.Background(), mockFolder, mockRequester.Username)
	assert.NoError(t, err, "Failed to create folder entry")

	mockFile = file_schema.NewUserFileBuilder().
		WithFilename("testfile").
		WithExtension("txt").
		WithSize(1024).
		WithContentType("text/plain").
		WithPath(mockRequester.Username).
		WithOwnerID(mockRequester.ID)

	mockFile.ID, err = connection.CreateFileEntry(context.Background(), mockFile, mockRequester)
	assert.NoError(t, err, "Failed to create file entry")

	return func() {
		cleanupMkdir()
		cleanupRename()
		cleanupRemoveAll()
		cleanupCreate()
		cleanupRemove()
		cleanupCopy()
		cleanupStat()
	}
}

func cleanupFileTestEnvironment(t *testing.T, cleanupFunc cleanupFunc) {
	for _, table := range setup.RequiredDBTables {
		err := connection.DropTable(context.Background(), table.Name)
		if err != nil {
			t.Errorf("Failed to drop table %s: %v", table.Name, err)
		}
	}
	if cleanupFunc != nil {
		cleanupFunc()
	}
}

func TestCreateFileEntry(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	fileID, err := con_copy.CreateFileEntry(context.Background(), mockFile, mockRequester)
	assert.NoError(t, err, "Failed to create file entry")
	assert.Greater(t, fileID, 0, "File ID should be greater than 0")
}

func TestExtractUserFiles(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	files, err := con_copy.ExtractUserFiles(context.Background(), mockRequester.ID)
	assert.NoError(t, err, "Failed to extract user files")
	assert.NotEmpty(t, files, "Files should not be empty")
	assert.Equal(t, "testfile", files[0].Filename, "Filename should match")
}

func TestRenameFile(t *testing.T) {
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

	newPath := "renamedfile.txt"

	renamedPath, err := con_copy.RenameFile(context.Background(), oldFile, newPath, mockRequester)
	assert.NoError(t, err, "Failed to rename file")
	assert.Equal(t, mockRequester.Username+"/"+newPath, renamedPath, "Renamed path should match")
}

func TestUpdateFileEntry(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	file := &file_schema.UserFile{
		ID:          1,
		Filename:    "updatedfile",
		Extension:   "txt",
		Path:        "testuser",
		Size:        2048,
		ContentType: "text/plain",
	}

	err := con_copy.UpdateFileEntry(context.Background(), file, mockRequester.Username)
	assert.NoError(t, err, "Failed to update file entry")
}

func TestDeleteFile(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	deletedID, err := con_copy.DeleteFile(context.Background(), mockFile, mockRequester.Username)
	assert.NoError(t, err, "Failed to delete file entry")
	assert.Equal(t, mockFile.ID, deletedID, "Deleted file ID should match")
}

func TestGetFile(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	mockFile.Filename = "renamedfile"

	file, err := con_copy.GetFile(context.Background(), mockFile, mockRequester)
	assert.NoError(t, err, "Failed to retrieve file entry")
	assert.Equal(t, mockFile.Filename, file.Filename, "Filename should match")
}

func TestExtractAllFiles(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	files, err := con_copy.ExtractAllFiles(context.Background())
	assert.NoError(t, err, "Failed to extract all files")
	assert.NotEmpty(t, files, "Files should not be empty")
}

func TestExtractFilesFromDirectory(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	directory := "testuser"
	files, err := con_copy.ExtractFilesFromDirectory(context.Background(), directory, mockRequester)
	assert.NoError(t, err, "Failed to extract files from directory")
	assert.NotEmpty(t, files, "Files should not be empty")
	assert.Equal(t, mockFile.Filename, files[0].Filename, "Filename should match")
}

func TestSaveFile(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	fileHeader := file_schema.NewUserFileBuilder().
		WithFilename("newfile").
		WithExtension("txt").
		WithPath("testuser").
		WithSize(512).
		WithContentType("text/plain").
		WithOwnerID(mockRequester.ID)

	mockFile := new(multipart.File)
	fileID, savedFileName, err := con_copy.SaveFile(context.Background(), fileHeader, mockFile, mockRequester)
	assert.NoError(t, err, "Failed to save file to filesystem")
	assert.Greater(t, fileID, 0, "File ID should be greater than 0")
	assert.Equal(t, "newfile.txt", savedFileName, "Saved file name should match")
}

func TestEnsureFolderExists(t *testing.T) {
	cleanup := setupFileTestEnvironment(t)
	defer cleanupFileTestEnvironment(t, cleanup)

	con_copy := connection.Copy()
	EndTransaction := con_copy.UseTransaction()
	assert.NotNil(t, EndTransaction)
	defer EndTransaction()

	relativePath := []string{"testuser"}
	err := con_copy.EnsureFolderExists(context.Background(), relativePath, mockRequester)
	assert.NoError(t, err, "Failed to ensure folder exists when relative path is a single segment")

	relativePath = []string{"testuser", "newfolder"}
	err = con_copy.EnsureFolderExists(context.Background(), relativePath, mockRequester)
	assert.NoError(t, err, "Failed to ensure folder exists when relative path has multiple segments")
}

func TestPreventDuplicateFileName(t *testing.T) {
	os.Setenv("STORAGE_PATH", "storage/") // Set environment variable for testing

	mockOS := new(mocks.MockOS)

	cleanupStat := test.MockFunction(&file.Stat, mockOS.Stat, func() {})
	defer cleanupStat()

	// Test cases
	t.Run("Returns original path if file does not exist", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/file.txt").Return(nil, os.ErrNotExist).Once() // Allow two calls

		path := "path/to/file.txt"
		newPath, err := file.PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.Equal(t, "storage/"+path, newPath, "Path should match the original path")
	})

	t.Run("Generates new path if file exists", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/file.txt").Return(nil, nil).Once()              // Allow two calls
		mockOS.On("Stat", "storage/path/to/file_1.txt").Return(nil, os.ErrNotExist).Once() // File exists

		path := "path/to/file.txt"
		newPath, err := file.PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.NotEqual(t, path, newPath, "Path should be updated to avoid duplication")
		assert.Equal(t, "storage/path/to/file_1.txt", newPath, "Generated path should match expected value")
	})

	t.Run("Handles multiple duplicates correctly", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/file.txt").Return(nil, nil).Once()              // File exists
		mockOS.On("Stat", "storage/path/to/file_1.txt").Return(nil, nil).Once()            // File exists
		mockOS.On("Stat", "storage/path/to/file_2.txt").Return(nil, os.ErrNotExist).Once() // File does not exist

		path := "path/to/file.txt"
		newPath, err := file.PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.Equal(t, "storage/path/to/file_2.txt", newPath, "Generated path should match expected value after handling multiple duplicates")
	})

	t.Run("Panics for invalid file path", func(t *testing.T) {
		assert.PanicsWithValue(t, "Invalid file path: ", func() {
			file.PreventDuplicateFileName("")
		}, "PreventDuplicateFileName should panic for invalid file path")
	})

	t.Run("Handles paths with special characters", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/special@file.txt").Return(nil, os.ErrNotExist).Once() // File does not exist

		path := "path/to/special@file.txt"
		newPath, err := file.PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.Equal(t, "storage/"+path, newPath, "Path should match the original path")
	})

	t.Run("Handles paths with no extension", func(t *testing.T) {
		path := "path/to/file"
		assert.Panics(t, func() {
			file.PreventDuplicateFileName(path)
		}, "PreventDuplicateFileName should panic for invalid file path")
	})
}
