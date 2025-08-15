package fschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBuilderSuite(t *testing.T) {
	t.Run("TestNewUserFileBuilder", func(t *testing.T) {
		newFile := NewUserFileBuilder()

		assert.NotNil(t, newFile, "New file should not be nil")
		assert.Equal(t, defaultFile, *newFile, "New file should match the default file")
	})

	t.Run("TestWithID", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithID(1)
		assert.Equal(t, 1, newFile.ID, "File ID should be set to 1")
	})

	t.Run("TestWithIDPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithID(0)
		}, "Setting file ID to 0 should panic")
	})

	t.Run("TestWithFilename", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithFilename("example.txt")
		assert.Equal(t, "example.txt", newFile.Filename, "Filename should be set to 'example.txt'")
	})

	t.Run("TestWithFilenamePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithFilename("")
		}, "Setting an empty filename should panic")
	})

	t.Run("TestWithExtension", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithExtension("txt")
		assert.Equal(t, "txt", newFile.Extension, "Extension should be set to 'txt'")
	})

	t.Run("TestWithExtensionPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithExtension("")
		}, "Setting an empty extension should panic")

		assert.Panics(t, func() {
			NewUserFileBuilder().WithExtension(".txt")
		}, "Setting an extension starting with a dot should panic")
	})

	t.Run("TestWithPath", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithPath("john_doe/documents")
		assert.Equal(t, "john_doe/documents", newFile.Path, "Path should be set to 'john_doe/documents'")
	})

	t.Run("TestWithPathPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithPath("")
		}, "Setting an empty path should panic")

		assert.Panics(t, func() {
			NewUserFileBuilder().WithPath("/invalid_path/")
		}, "Setting a path starting or ending with a slash should panic")
	})

	t.Run("TestWithDirectory", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithDirectory(42)
		assert.Equal(t, 42, newFile.Directory, "Directory ID should be set to 42")
	})

	t.Run("TestWithDirectoryPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithDirectory(0)
		}, "Setting directory ID to 0 should panic")
	})

	t.Run("TestWithSize", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithSize(1024)
		assert.Equal(t, int64(1024), newFile.Size, "File size should be set to 1024")
	})

	t.Run("TestWithSizePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithSize(0)
		}, "Setting file size to 0 should panic")
	})

	t.Run("TestWithContentType", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithContentType("text/plain")
		assert.Equal(t, "text/plain", newFile.ContentType, "Content type should be set to 'text/plain'")
	})

	t.Run("TestWithContentTypePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithContentType("")
		}, "Setting an empty content type should panic")
	})

	t.Run("TestWithOwnerID", func(t *testing.T) {
		newFile := NewUserFileBuilder().WithOwnerID(100)
		assert.Equal(t, 100, newFile.OwnerID, "Owner ID should be set to 100")
	})

	t.Run("TestWithOwnerIDPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithOwnerID(0)
		}, "Setting owner ID to 0 should panic")
	})

	t.Run("TestWithMetadata", func(t *testing.T) {
		metadata := "file metadata"
		newFile := NewUserFileBuilder().WithMetadata(&metadata)
		assert.Equal(t, &metadata, newFile.Metadata, "Metadata should be set to 'file metadata'")
	})

	t.Run("TestWithMetadataPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFileBuilder().WithMetadata(nil)
		}, "Setting metadata to nil should panic")

		emptyMetadata := ""
		assert.Panics(t, func() {
			NewUserFileBuilder().WithMetadata(&emptyMetadata)
		}, "Setting metadata to an empty string should panic")
	})
}
