package db

import (
	"context"
	"homeserver/pkg/db/file"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"homeserver/test"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileMethods(t *testing.T) {
	mockOS := new(MockOS)
	mockOS.On("MkdirAll", mock.Anything, mock.Anything).Return(nil)
	mockOS.On("Rename", mock.Anything, mock.Anything).Return(nil)
	mockOS.On("RemoveAll", mock.Anything).Return(nil)
	mockOS.On("Create", mock.Anything).Return(nil, nil)
	mockOS.On("Remove", mock.Anything).Return(nil)

	cleanupMkdir := test.MockFunction(&file.MkdirAll, mockOS.MkdirAll, func() {})
	defer cleanupMkdir()
	cleanupRename := test.MockFunction(&file.Rename, mockOS.Rename, func() {})
	defer cleanupRename()
	cleanupRemoveAll := test.MockFunction(&file.RemoveAll, mockOS.RemoveAll, func() {})
	defer cleanupRemoveAll()
	cleanupCreate := test.MockFunction(&file.Create, mockOS.Create, func() {})
	defer cleanupCreate()
	cleanupRemove := test.MockFunction(&file.Remove, mockOS.Remove, func() {})
	defer cleanupRemove()

	mockIO := new(MockIO)
	mockIO.On("Copy", mock.Anything, mock.Anything).Return(int64(0), nil)

	cleanupCopy := test.MockFunction(&file.Copy, mockIO.Copy, func() {})
	defer cleanupCopy()

	assert.NotNil(t, connection, "Database connection should not be nil")

	err := setup.DatabaseSetup(connection)
	assert.NoError(t, err, "Failed to check database tables")

	mockRequester := user_schema.NewUserBuilder().
		WithUsername("testuser").
		WithEmail("testuser@example.com").
		WithPassword("password123")

	mockRequester, err = connection.CreateUserEntry(context.Background(), mockRequester)
	assert.NoError(t, err, "Failed to create user entry")

	mockFolder := folder_schema.NewUserFolderBuilder().
		WithName(mockRequester.Username).
		WithOwnerID(mockRequester.ID)

	err = connection.CreateRootFolderEntry(context.Background(), mockFolder, mockRequester.Username)
	assert.NoError(t, err, "Failed to create folder entry")

	mockFile := file_schema.NewUserFileBuilder().
		WithFilename("testfile").
		WithExtension("txt").
		WithSize(1024).
		WithContentType("text/plain").
		WithPath(mockRequester.Username).
		WithOwnerID(mockRequester.ID)

	mockFile.ID, err = connection.CreateFileEntry(context.Background(), mockFile, mockRequester)
	assert.NoError(t, err, "Failed to create file entry")

	defer func() {
		for _, table := range setup.RequiredDBTables {
			err := connection.DropTable(context.Background(), table.Name)
			assert.NoError(t, err, "Failed to drop table %s", table.Name)
		}
	}()

	t.Run("CreateFileEntry inserts a file entry", func(t *testing.T) {
		assert.NotNil(t, connection, "Connection should not be nil")

		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		file := &file_schema.UserFile{
			Filename:    "testfile",
			Extension:   "txt",
			Path:        "testuser",
			Size:        1024,
			ContentType: "text/plain",
			OwnerID:     mockRequester.ID,
		}

		fileID, err := connection.CreateFileEntry(context.Background(), file, mockRequester)
		assert.NoError(t, err, "Failed to create file entry")
		assert.Greater(t, fileID, 0, "File ID should be greater than 0")
	})

	t.Run("ExtractUserFiles retrieves user files", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		files, err := connection.ExtractUserFiles(context.Background(), mockRequester.ID)
		assert.NoError(t, err, "Failed to extract user files")
		assert.NotEmpty(t, files, "Files should not be empty")
		assert.Equal(t, "testfile", files[0].Filename, "Filename should match")
	})

	t.Run("RenameFile renames a file", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		oldPath := "testfile.txt"
		newPath := "renamedfile.txt"

		renamedPath, err := connection.RenameFile(context.Background(), oldPath, newPath, mockRequester)
		assert.NoError(t, err, "Failed to rename file")
		assert.Equal(t, mockRequester.Username+"/"+newPath, renamedPath, "Renamed path should match")
	})

	t.Run("UpdateFileEntry updates file metadata", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
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

		err := connection.UpdateFileEntry(context.Background(), file)
		assert.NoError(t, err, "Failed to update file entry")
	})

	t.Run("DeleteFile removes a file entry", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		file := &file_schema.UserFile{
			ID:        1,
			Filename:  "updatedfile",
			Extension: "txt",
			OwnerID:   mockRequester.ID,
			Path:      mockRequester.Username,
		}

		deletedID, err := connection.DeleteFile(context.Background(), file, mockRequester.Username)
		assert.NoError(t, err, "Failed to delete file entry")
		assert.Equal(t, file.ID, deletedID, "Deleted file ID should match")
	})

	t.Run("GetFile retrieves a file entry", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		mockFile.Filename = "renamedfile"

		file, err := connection.GetFile(context.Background(), mockFile, mockRequester)
		assert.NoError(t, err, "Failed to retrieve file entry")
		assert.Equal(t, mockFile.Filename, file.Filename, "Filename should match")
	})

	t.Run("ExtractAllFiles retrieves all files independently of users", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		files, err := connection.ExtractAllFiles(context.Background())
		assert.NoError(t, err, "Failed to extract all files")
		assert.NotEmpty(t, files, "Files should not be empty")
	})

	t.Run("ExtractFilesFromDirectory retrieves files from a specific directory", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		directory := "testuser"
		files, err := connection.ExtractFilesFromDirectory(context.Background(), directory, mockRequester)
		assert.NoError(t, err, "Failed to extract files from directory")
		assert.NotEmpty(t, files, "Files should not be empty")
		assert.Equal(t, mockFile.Filename, files[0].Filename, "Filename should match")
	})

	t.Run("SaveFile saves a file to the filesystem", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		fileHeader := &file_schema.UserFile{
			Filename:    "newfile",
			Extension:   "txt",
			Path:        "testuser",
			Size:        512,
			ContentType: "text/plain",
			OwnerID:     mockRequester.ID,
		}

		mockFile := new(multipart.File)
		fileID, savedFileName, err := connection.SaveFile(context.Background(), fileHeader, mockFile, mockRequester)
		assert.NoError(t, err, "Failed to save file to filesystem")
		assert.Greater(t, fileID, 0, "File ID should be greater than 0")
		assert.Equal(t, "newfile.txt", savedFileName, "Saved file name should match")
	})

	t.Run("EnsureFolderExists checks if a folder exists and creates it if not", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		relativePath := []string{"testuser"}
		err := connection.EnsureFolderExists(context.Background(), relativePath, mockRequester)
		assert.NoError(t, err, "Failed to ensure folder exists when relative path is a single segment")

		relativePath = []string{"testuser", "newfolder"}
		err = connection.EnsureFolderExists(context.Background(), relativePath, mockRequester)
		assert.NoError(t, err, "Failed to ensure folder exists when relative path has multiple segments")
	})
}
