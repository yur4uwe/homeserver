package fschema

import (
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks/mocksos"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBuilderSuite(t *testing.T) {
	t.Run("TestNewUserFileBuilder", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder()

		assert.NotNil(t, newFile, "New file should not be nil")
		assert.Equal(t, fschema.DefaultFile, *newFile, "New file should match the default file")
	})

	t.Run("TestWithID", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithID(1)
		assert.Equal(t, 1, newFile.ID, "File ID should be set to 1")
	})

	t.Run("TestWithIDPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithID(0)
		}, "Setting file ID to 0 should panic")
	})

	t.Run("TestWithFilename", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithFilename("example.txt")
		assert.Equal(t, "example.txt", newFile.Filename, "Filename should be set to 'example.txt'")
	})

	t.Run("TestWithFilenamePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithFilename("")
		}, "Setting an empty filename should panic")
	})

	t.Run("TestWithExtension", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithExtension("txt")
		assert.Equal(t, "txt", newFile.Extension, "Extension should be set to 'txt'")
	})

	t.Run("TestWithExtensionPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithExtension("")
		}, "Setting an empty extension should panic")

		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithExtension(".txt")
		}, "Setting an extension starting with a dot should panic")
	})

	t.Run("TestWithPath", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithPath("john_doe/documents")
		assert.Equal(t, "john_doe/documents", newFile.Path, "Path should be set to 'john_doe/documents'")
	})

	t.Run("TestWithPathPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithPath("")
		}, "Setting an empty path should panic")

		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithPath("/invalid_path/")
		}, "Setting a path starting or ending with a slash should panic")
	})

	t.Run("TestWithDirectory", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithDirectory(42)
		assert.Equal(t, 42, newFile.Directory, "Directory ID should be set to 42")
	})

	t.Run("TestWithDirectoryPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithDirectory(0)
		}, "Setting directory ID to 0 should panic")
	})

	t.Run("TestWithSize", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithSize(1024)
		assert.Equal(t, int64(1024), newFile.Size, "File size should be set to 1024")
	})

	t.Run("TestWithSizePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithSize(0)
		}, "Setting file size to 0 should panic")
	})

	t.Run("TestWithContentType", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithContentType("text/plain")
		assert.Equal(t, "text/plain", newFile.ContentType, "Content type should be set to 'text/plain'")
	})

	t.Run("TestWithContentTypePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithContentType("")
		}, "Setting an empty content type should panic")
	})

	t.Run("TestWithOwnerID", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithOwnerID(100)
		assert.Equal(t, 100, newFile.OwnerID, "Owner ID should be set to 100")
	})

	t.Run("TestWithOwnerIDPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithOwnerID(0)
		}, "Setting owner ID to 0 should panic")
	})

	t.Run("TestWithMetadata", func(t *testing.T) {
		metadata := "file metadata"
		newFile := fschema.NewUserFileBuilder().WithMetadata(&metadata)
		assert.Equal(t, &metadata, newFile.Metadata, "Metadata should be set to 'file metadata'")
	})

	t.Run("TestWithMetadataPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithMetadata(nil)
		}, "Setting metadata to nil should panic")

		emptyMetadata := ""
		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithMetadata(&emptyMetadata)
		}, "Setting metadata to an empty string should panic")
	})

	t.Run("TestWithFullName", func(t *testing.T) {
		newFile := fschema.NewUserFileBuilder().WithFullName("archive.tar.gz")
		assert.Equal(t, "archive.tar", newFile.Filename, "Filename should be 'archive.tar'")
		assert.Equal(t, "gz", newFile.Extension, "Extension should be 'gz'")

		newFile2 := fschema.NewUserFileBuilder().WithFullName("file.txt")
		assert.Equal(t, "file", newFile2.Filename, "Filename should be 'file'")
		assert.Equal(t, "txt", newFile2.Extension, "Extension should be 'txt'")

		assert.Panics(t, func() {
			fschema.NewUserFileBuilder().WithFullName("")
		}, "Setting an empty full name should panic")
	})

	t.Run("TestAbsolutePath", func(t *testing.T) {
		mockOS := &mocksos.MockOS{}
		mockOS.On("Getenv", "STORAGE_PATH").Return("/storage/")
		helpers.MockFunction(&fschema.Getenv, mockOS.Getenv, func() {})

		newFile := fschema.NewUserFileBuilder().
			WithFilename("report").
			WithExtension("pdf").
			WithPath("john/reports")
		expected := "/storage/john/reports/report.pdf"
		assert.Equal(t, expected, newFile.AbsolutePath())

		newFile2 := fschema.NewUserFileBuilder().
			WithFilename("file").
			WithExtension("txt")
		expected2 := "/storage/file.txt"
		assert.Equal(t, expected2, newFile2.AbsolutePath())
	})

	t.Run("TestFieldsList", func(t *testing.T) {
		fields := fschema.FieldsList()
		assert.Contains(t, fields, "id")
		assert.Contains(t, fields, "filename")
		assert.Contains(t, fields, "extension")
		assert.Contains(t, fields, "path")
		assert.Contains(t, fields, "directory_id")
		assert.Contains(t, fields, "size")
		assert.Contains(t, fields, "content_type")
		assert.Contains(t, fields, "owner_id")
		assert.Contains(t, fields, "created_at")
		assert.Contains(t, fields, "updated_at")
		assert.Contains(t, fields, "metadata")

		withoutFields := fschema.FieldsList("id", "metadata")
		assert.NotContains(t, withoutFields, "id")
		assert.NotContains(t, withoutFields, "metadata")
	})
}
