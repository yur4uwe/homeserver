package fschema

import (
	"fmt"
	"homeserver/constants/fconsts"
	fschema "homeserver/db/repositories/file_repo/schema"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileChecksSuite(t *testing.T) {
	t.Run("TestCheckFileID", func(t *testing.T) {
		file := &fschema.UserFile{ID: 1}
		err := fschema.CheckFileID(file)
		assert.NoError(t, err, "CheckFileID should not return an error for valid ID")

		file = &fschema.UserFile{ID: 0}
		err = fschema.CheckFileID(file)
		assert.Error(t, err, "CheckFileID should return an error for invalid ID")
		assert.EqualError(t, err, fconsts.InvalidIDError, "Error message should match")
	})

	t.Run("TestCheckFilename", func(t *testing.T) {
		file := &fschema.UserFile{Filename: "example.txt"}
		err := fschema.CheckFilename(file)
		assert.NoError(t, err, "CheckFilename should not return an error for valid filename")

		file = &fschema.UserFile{Filename: ""}
		err = fschema.CheckFilename(file)
		assert.Error(t, err, "CheckFilename should return an error for empty filename")
		assert.EqualError(t, err, fconsts.MissingNameError, "Error message should match")
	})

	t.Run("TestCheckExtension", func(t *testing.T) {
		file := &fschema.UserFile{Extension: "txt"}
		err := fschema.CheckExtension(file)
		assert.NoError(t, err, "CheckExtension should not return an error for valid extension")

		file = &fschema.UserFile{Extension: ""}
		err = fschema.CheckExtension(file)
		assert.Error(t, err, "CheckExtension should return an error for empty extension")
		assert.EqualError(t, err, fconsts.InvalidExtensionError, "Error message should match")

		file = &fschema.UserFile{Extension: ".txt"}
		err = fschema.CheckExtension(file)
		assert.Error(t, err, "CheckExtension should return an error for extension starting with a dot")
		assert.EqualError(t, err, fconsts.InvalidExtensionError, "Error message should match")
	})

	t.Run("TestCheckPath", func(t *testing.T) {
		checkPath := fschema.CheckPath("john_doe")

		file := &fschema.UserFile{Path: "john_doe/documents"}
		err := checkPath(file)
		assert.NoError(t, err, "CheckPath should not return an error for valid path")

		file = &fschema.UserFile{Path: ""}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for empty path")
		assert.EqualError(t, err, fconsts.MissingPathError, "Error message should match")

		file = &fschema.UserFile{Path: "/invalid_path"}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for path starting with a slash")
		assert.EqualError(t, err, fconsts.InvalidPathStartError, "Error message should match")

		file = &fschema.UserFile{Path: "invalid_path/"}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for path ending with a slash")
		assert.EqualError(t, err, fconsts.InvalidPathEndError, "Error message should match")

		file = &fschema.UserFile{Path: "other_user/documents"}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for path not starting with owner's name")
		assert.EqualError(t, err, fconsts.MissingOwnerFolderError, "Error message should match")
	})

	t.Run("TestCheckDirectory", func(t *testing.T) {
		file := &fschema.UserFile{Directory: 42}
		err := fschema.CheckDirectory(file)
		assert.NoError(t, err, "CheckDirectory should not return an error for valid directory ID")

		file = &fschema.UserFile{Directory: 0}
		err = fschema.CheckDirectory(file)
		assert.Error(t, err, "CheckDirectory should return an error for invalid directory ID")
		assert.EqualError(t, err, fconsts.InvalidDirectoryError, "Error message should match")
	})

	t.Run("TestCheckSize", func(t *testing.T) {
		file := &fschema.UserFile{Size: 1024}
		err := fschema.CheckSize(file)
		assert.NoError(t, err, "CheckSize should not return an error for valid size")

		file = &fschema.UserFile{Size: 0}
		err = fschema.CheckSize(file)
		assert.Error(t, err, "CheckSize should return an error for invalid size")
		assert.EqualError(t, err, fconsts.InvalidSizeError, "Error message should match")
	})

	t.Run("TestCheckContentType", func(t *testing.T) {
		file := &fschema.UserFile{ContentType: "text/plain"}
		err := fschema.CheckContentType(file)
		assert.NoError(t, err, "CheckContentType should not return an error for valid content type")

		file = &fschema.UserFile{ContentType: ""}
		err = fschema.CheckContentType(file)
		assert.Error(t, err, "CheckContentType should return an error for empty content type")
		assert.EqualError(t, err, fconsts.InvalidContentTypeError, "Error message should match")
	})

	t.Run("TestCheckOwnerID", func(t *testing.T) {
		file := &fschema.UserFile{OwnerID: 100}
		err := fschema.CheckOwnerID(file)
		assert.NoError(t, err, "CheckOwnerID should not return an error for valid owner ID")

		file = &fschema.UserFile{OwnerID: 0}
		err = fschema.CheckOwnerID(file)
		assert.Error(t, err, "CheckOwnerID should return an error for invalid owner ID")
		assert.EqualError(t, err, fconsts.InvalidOwnerIDError, "Error message should match")
	})

	t.Run("TestCheckMetadata", func(t *testing.T) {
		metadata := "file metadata"
		file := &fschema.UserFile{Metadata: &metadata}
		err := fschema.CheckMetadata(file)
		assert.NoError(t, err, "CheckMetadata should not return an error for valid metadata")

		file = &fschema.UserFile{Metadata: nil}
		err = fschema.CheckMetadata(file)
		assert.Error(t, err, "CheckMetadata should return an error for nil metadata")
		assert.EqualError(t, err, fconsts.MissingMetadataError, "Error message should match")

		emptyMetadata := ""
		file = &fschema.UserFile{Metadata: &emptyMetadata}
		err = fschema.CheckMetadata(file)
		assert.Error(t, err, "CheckMetadata should return an error for empty metadata")
		assert.EqualError(t, err, fconsts.MissingMetadataError, "Error message should match")
	})

	t.Run("TestShouldHaveFields", func(t *testing.T) {
		metadata := "file metadata"

		file := &fschema.UserFile{
			ID:          0,
			Filename:    "",
			Extension:   ".txt",
			Path:        "/invalid_path/",
			Directory:   0,
			Size:        -1,
			ContentType: "",
			OwnerID:     0,
			Metadata:    nil,
		}

		expectedErr := "file validation failed:\n" +
			fconsts.ID + ": " + fconsts.InvalidIDError + "\n" +
			fconsts.Filename + ": " + fconsts.MissingNameError + "\n" +
			fconsts.Extension + ": " + fconsts.InvalidExtensionError + "\n" +
			fconsts.Path + ": " + fconsts.InvalidPathStartError + "\n" +
			fconsts.Directory + ": " + fconsts.InvalidDirectoryError + "\n" +
			fconsts.Size + ": " + fconsts.InvalidSizeError + "\n" +
			fconsts.ContentType + ": " + fconsts.InvalidContentTypeError + "\n" +
			fconsts.OwnerID + ": " + fconsts.InvalidOwnerIDError + "\n" +
			fconsts.Metadata + ": " + fconsts.MissingMetadataError

		err := file.ShouldHaveFields("john_doe")(
			fconsts.ID, fconsts.Filename, fconsts.Extension, fconsts.Path,
			fconsts.Directory, fconsts.Size, fconsts.ContentType, fconsts.OwnerID, fconsts.Metadata)
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid fields")
		assert.EqualError(t, err, expectedErr, "Error message should match")

		file = &fschema.UserFile{
			ID:          1,
			Filename:    "example.txt",
			Extension:   "txt",
			Path:        "john_doe/documents",
			Directory:   42,
			Size:        1024,
			ContentType: "text/plain",
			OwnerID:     100,
			Metadata:    &metadata,
		}

		err = file.ShouldHaveFields("john_doe")(
			fconsts.ID, fconsts.Filename, fconsts.Extension, fconsts.Path,
			fconsts.Directory, fconsts.Size, fconsts.ContentType, fconsts.OwnerID, fconsts.Metadata)
		assert.NoError(t, err, "ShouldHaveFields should not return an error for valid fields")

		err = file.ShouldHaveFields("john_doe")()
		assert.Error(t, err, "ShouldHaveFields should return an error when no fields are specified")
		assert.EqualError(t, err, fconsts.NoFieldsSpecifiedError, "Error message should match")

		err = file.ShouldHaveFields("john_doe")("invalid_field")
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid field")
		assert.EqualError(t, err, "file validation failed:\n"+fmt.Sprintf(fconsts.InvalidFileFieldError, "invalid_field"), "Error message should match")
	})
}
