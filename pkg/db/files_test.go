package db

import (
	"context"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/setup"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileMethods(t *testing.T) {
	assert.NotNil(t, connection, "Database connection should not be nil")

	// Create the files table
	err := setup.CheckDatabaseTables(connection)
	assert.NoError(t, err, "Failed to check database tables")
	defer func() {
		for _, table := range setup.RequiredDBTables {
			if err := connection.DropTable(context.Background(), table.Name); err != nil {
				t.Errorf("Failed to drop table %s: %v", table.Name, err)
			}
		}
	}()

	// Mock user for testing
	requester := &user_schema.User{
		ID:       1,
		Username: "testuser",
	}

	t.Run("CreateFileEntry inserts a file entry", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		file := &file_schema.UserFile{
			Filename:    "testfile",
			Extension:   "txt",
			Path:        "testuser/files",
			Size:        1024,
			ContentType: "text/plain",
		}

		fileID, err := connection.CreateFileEntry(context.Background(), file, requester)
		assert.NoError(t, err, "Failed to create file entry")
		assert.Greater(t, fileID, 0, "File ID should be greater than 0")
	})

	t.Run("ExtractUserFiles retrieves user files", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		files, err := connection.ExtractUserFiles(context.Background(), requester.ID)
		assert.NoError(t, err, "Failed to extract user files")
		assert.NotEmpty(t, files, "Files should not be empty")
		assert.Equal(t, "testfile", files[0].Filename, "Filename should match")
	})

	t.Run("RenameFile renames a file", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		oldPath := "testuser/files/testfile.txt"
		newPath := "testuser/files/renamedfile.txt"

		renamedPath, err := connection.RenameFile(context.Background(), oldPath, newPath, requester)
		assert.NoError(t, err, "Failed to rename file")
		assert.Equal(t, newPath, renamedPath, "Renamed path should match")
	})

	t.Run("UpdateFileEntry updates file metadata", func(t *testing.T) {
		EndTransaction := connection.UseTransaction()
		assert.NotNil(t, EndTransaction)
		defer EndTransaction()

		file := &file_schema.UserFile{
			ID:          1,
			Filename:    "updatedfile",
			Extension:   "txt",
			Path:        "testuser/files",
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
			Filename: "updatedfile",
			OwnerID:  requester.ID,
		}

		deletedID, err := connection.DeleteFile(context.Background(), file, requester.Username)
		assert.NoError(t, err, "Failed to delete file entry")
		assert.Equal(t, file.ID, deletedID, "Deleted file ID should match")
	})
}
