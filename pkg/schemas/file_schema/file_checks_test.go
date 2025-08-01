package file_schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileChecksSuite(t *testing.T) {
	t.Run("TestCheckFileID", func(t *testing.T) {
		file := &UserFile{ID: 1}
		result, err := CheckFileID(file, nil)
		assert.NoError(t, err, "CheckFileID should not return an error for valid ID")
		assert.Equal(t, file, result, "CheckFileID should return the same file object")

		file = &UserFile{ID: 0}
		_, err = CheckFileID(file, nil)
		assert.Error(t, err, "CheckFileID should return an error for invalid ID")
		assert.EqualError(t, err, "file ID is invalid: ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckFilename", func(t *testing.T) {
		file := &UserFile{Filename: "example.txt"}
		result, err := CheckFilename(file, nil)
		assert.NoError(t, err, "CheckFilename should not return an error for valid filename")
		assert.Equal(t, file, result, "CheckFilename should return the same file object")

		file = &UserFile{Filename: ""}
		_, err = CheckFilename(file, nil)
		assert.Error(t, err, "CheckFilename should return an error for empty filename")
		assert.EqualError(t, err, "filename is invalid: filename cannot be empty", "Error message should match")
	})

	t.Run("TestCheckExtension", func(t *testing.T) {
		file := &UserFile{Extension: "txt"}
		result, err := CheckExtension(file, nil)
		assert.NoError(t, err, "CheckExtension should not return an error for valid extension")
		assert.Equal(t, file, result, "CheckExtension should return the same file object")

		file = &UserFile{Extension: ""}
		_, err = CheckExtension(file, nil)
		assert.Error(t, err, "CheckExtension should return an error for empty extension")
		assert.EqualError(t, err, "extension is invalid: extension cannot be empty", "Error message should match")

		file = &UserFile{Extension: ".txt"}
		_, err = CheckExtension(file, nil)
		assert.Error(t, err, "CheckExtension should return an error for extension starting with a dot")
		assert.EqualError(t, err, "extension is invalid: extension cannot start with a dot", "Error message should match")
	})

	t.Run("TestCheckPath", func(t *testing.T) {
		checkPath := CheckPath("john_doe")

		file := &UserFile{Path: "john_doe/documents"}
		result, err := checkPath(file, nil)
		assert.NoError(t, err, "CheckPath should not return an error for valid path")
		assert.Equal(t, file, result, "CheckPath should return the same file object")

		file = &UserFile{Path: ""}
		_, err = checkPath(file, nil)
		assert.Error(t, err, "CheckPath should return an error for empty path")
		assert.EqualError(t, err, "path is invalid: path cannot be empty", "Error message should match")

		file = &UserFile{Path: "/invalid_path"}
		_, err = checkPath(file, nil)
		assert.Error(t, err, "CheckPath should return an error for path starting with a slash")
		assert.EqualError(t, err, "path is invalid: path cannot start with a slash", "Error message should match")

		file = &UserFile{Path: "invalid_path/"}
		_, err = checkPath(file, nil)
		assert.Error(t, err, "CheckPath should return an error for path ending with a slash")
		assert.EqualError(t, err, "path is invalid: path cannot end with a slash", "Error message should match")

		file = &UserFile{Path: "other_user/documents"}
		_, err = checkPath(file, nil)
		assert.Error(t, err, "CheckPath should return an error for path not starting with owner's name")
		assert.EqualError(t, err, "path is invalid: path must start with the owner's name 'john_doe'", "Error message should match")
	})

	t.Run("TestCheckDirectory", func(t *testing.T) {
		file := &UserFile{Directory: 42}
		result, err := CheckDirectory(file, nil)
		assert.NoError(t, err, "CheckDirectory should not return an error for valid directory ID")
		assert.Equal(t, file, result, "CheckDirectory should return the same file object")

		file = &UserFile{Directory: 0}
		_, err = CheckDirectory(file, nil)
		assert.Error(t, err, "CheckDirectory should return an error for invalid directory ID")
		assert.EqualError(t, err, "directory is invalid: directory ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckSize", func(t *testing.T) {
		file := &UserFile{Size: 1024}
		result, err := CheckSize(file, nil)
		assert.NoError(t, err, "CheckSize should not return an error for valid size")
		assert.Equal(t, file, result, "CheckSize should return the same file object")

		file = &UserFile{Size: 0}
		_, err = CheckSize(file, nil)
		assert.Error(t, err, "CheckSize should return an error for invalid size")
		assert.EqualError(t, err, "size is invalid: file size must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckContentType", func(t *testing.T) {
		file := &UserFile{ContentType: "text/plain"}
		result, err := CheckContentType(file, nil)
		assert.NoError(t, err, "CheckContentType should not return an error for valid content type")
		assert.Equal(t, file, result, "CheckContentType should return the same file object")

		file = &UserFile{ContentType: ""}
		_, err = CheckContentType(file, nil)
		assert.Error(t, err, "CheckContentType should return an error for empty content type")
		assert.EqualError(t, err, "content type is invalid: content type cannot be empty", "Error message should match")
	})

	t.Run("TestCheckOwnerID", func(t *testing.T) {
		file := &UserFile{OwnerID: 100}
		result, err := CheckOwnerID(file, nil)
		assert.NoError(t, err, "CheckOwnerID should not return an error for valid owner ID")
		assert.Equal(t, file, result, "CheckOwnerID should return the same file object")

		file = &UserFile{OwnerID: 0}
		_, err = CheckOwnerID(file, nil)
		assert.Error(t, err, "CheckOwnerID should return an error for invalid owner ID")
		assert.EqualError(t, err, "owner ID is invalid: owner ID must be greater than 0", "Error message should match")
	})

	t.Run("TestCheckMetadata", func(t *testing.T) {
		metadata := "file metadata"
		file := &UserFile{Metadata: &metadata}
		result, err := CheckMetadata(file, nil)
		assert.NoError(t, err, "CheckMetadata should not return an error for valid metadata")
		assert.Equal(t, file, result, "CheckMetadata should return the same file object")

		file = &UserFile{Metadata: nil}
		_, err = CheckMetadata(file, nil)
		assert.Error(t, err, "CheckMetadata should return an error for nil metadata")
		assert.EqualError(t, err, "metadata is invalid: metadata cannot be nil", "Error message should match")

		emptyMetadata := ""
		file = &UserFile{Metadata: &emptyMetadata}
		_, err = CheckMetadata(file, nil)
		assert.Error(t, err, "CheckMetadata should return an error for empty metadata")
		assert.EqualError(t, err, "metadata is invalid: metadata cannot be empty", "Error message should match")
	})
}
