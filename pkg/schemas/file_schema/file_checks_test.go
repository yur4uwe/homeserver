package file_schema

import (
	"fmt"
	"homeserver/pkg/schemas"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileChecksSuite(t *testing.T) {
	t.Run("TestCheckFileID", func(t *testing.T) {
		file := &UserFile{ID: 1}
		err := CheckFileID(file)
		assert.NoError(t, err, "CheckFileID should not return an error for valid ID")

		file = &UserFile{ID: 0}
		err = CheckFileID(file)
		assert.Error(t, err, "CheckFileID should return an error for invalid ID")
		assert.EqualError(t, err, schemas.InvalidFileIDError, "Error message should match")
	})

	t.Run("TestCheckFilename", func(t *testing.T) {
		file := &UserFile{Filename: "example.txt"}
		err := CheckFilename(file)
		assert.NoError(t, err, "CheckFilename should not return an error for valid filename")

		file = &UserFile{Filename: ""}
		err = CheckFilename(file)
		assert.Error(t, err, "CheckFilename should return an error for empty filename")
		assert.EqualError(t, err, schemas.InvalidFilenameError, "Error message should match")
	})

	t.Run("TestCheckExtension", func(t *testing.T) {
		file := &UserFile{Extension: "txt"}
		err := CheckExtension(file)
		assert.NoError(t, err, "CheckExtension should not return an error for valid extension")

		file = &UserFile{Extension: ""}
		err = CheckExtension(file)
		assert.Error(t, err, "CheckExtension should return an error for empty extension")
		assert.EqualError(t, err, schemas.InvalidExtensionError, "Error message should match")

		file = &UserFile{Extension: ".txt"}
		err = CheckExtension(file)
		assert.Error(t, err, "CheckExtension should return an error for extension starting with a dot")
		assert.EqualError(t, err, schemas.InvalidExtensionError, "Error message should match")
	})

	t.Run("TestCheckPath", func(t *testing.T) {
		checkPath := CheckPath("john_doe")

		file := &UserFile{Path: "john_doe/documents"}
		err := checkPath(file)
		assert.NoError(t, err, "CheckPath should not return an error for valid path")

		file = &UserFile{Path: ""}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for empty path")
		assert.EqualError(t, err, schemas.InvalidFolderPathError, "Error message should match")

		file = &UserFile{Path: "/invalid_path"}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for path starting with a slash")
		assert.EqualError(t, err, schemas.InvalidFilePathError, "Error message should match")

		file = &UserFile{Path: "invalid_path/"}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for path ending with a slash")
		assert.EqualError(t, err, schemas.InvalidFilePathEndError, "Error message should match")

		file = &UserFile{Path: "other_user/documents"}
		err = checkPath(file)
		assert.Error(t, err, "CheckPath should return an error for path not starting with owner's name")
		assert.EqualError(t, err, fmt.Sprintf(schemas.InvalidFilePathOwnerError, "john_doe"), "Error message should match")
	})

	t.Run("TestCheckDirectory", func(t *testing.T) {
		file := &UserFile{Directory: 42}
		err := CheckDirectory(file)
		assert.NoError(t, err, "CheckDirectory should not return an error for valid directory ID")

		file = &UserFile{Directory: 0}
		err = CheckDirectory(file)
		assert.Error(t, err, "CheckDirectory should return an error for invalid directory ID")
		assert.EqualError(t, err, schemas.InvalidDirectoryError, "Error message should match")
	})

	t.Run("TestCheckSize", func(t *testing.T) {
		file := &UserFile{Size: 1024}
		err := CheckSize(file)
		assert.NoError(t, err, "CheckSize should not return an error for valid size")

		file = &UserFile{Size: 0}
		err = CheckSize(file)
		assert.Error(t, err, "CheckSize should return an error for invalid size")
		assert.EqualError(t, err, schemas.InvalidFileSizeError, "Error message should match")
	})

	t.Run("TestCheckContentType", func(t *testing.T) {
		file := &UserFile{ContentType: "text/plain"}
		err := CheckContentType(file)
		assert.NoError(t, err, "CheckContentType should not return an error for valid content type")

		file = &UserFile{ContentType: ""}
		err = CheckContentType(file)
		assert.Error(t, err, "CheckContentType should return an error for empty content type")
		assert.EqualError(t, err, schemas.InvalidContentTypeError, "Error message should match")
	})

	t.Run("TestCheckOwnerID", func(t *testing.T) {
		file := &UserFile{OwnerID: 100}
		err := CheckOwnerID(file)
		assert.NoError(t, err, "CheckOwnerID should not return an error for valid owner ID")

		file = &UserFile{OwnerID: 0}
		err = CheckOwnerID(file)
		assert.Error(t, err, "CheckOwnerID should return an error for invalid owner ID")
		assert.EqualError(t, err, schemas.InvalidFileOwnerIDError, "Error message should match")
	})

	t.Run("TestCheckMetadata", func(t *testing.T) {
		metadata := "file metadata"
		file := &UserFile{Metadata: &metadata}
		err := CheckMetadata(file)
		assert.NoError(t, err, "CheckMetadata should not return an error for valid metadata")

		file = &UserFile{Metadata: nil}
		err = CheckMetadata(file)
		assert.Error(t, err, "CheckMetadata should return an error for nil metadata")
		assert.EqualError(t, err, schemas.InvalidFileMetadataError, "Error message should match")

		emptyMetadata := ""
		file = &UserFile{Metadata: &emptyMetadata}
		err = CheckMetadata(file)
		assert.Error(t, err, "CheckMetadata should return an error for empty metadata")
		assert.EqualError(t, err, schemas.InvalidFileMetadataError, "Error message should match")
	})

	t.Run("TestShouldHaveFields", func(t *testing.T) {
		metadata := "file metadata"

		file := &UserFile{
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

		err := file.ShouldHaveFields("john_doe")(schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Directory, schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.Metadata)
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid fields")
		assert.EqualError(t, err, schemas.ValidationErrorPrefix+"FileID: "+schemas.InvalidFileIDError+"\nFilename: "+schemas.InvalidFilenameError+"\nExtension: "+schemas.InvalidExtensionError+"\nPath: "+schemas.InvalidFilePathError+"\nDirectory: "+schemas.InvalidDirectoryError+"\nSize: "+schemas.InvalidFileSizeError+"\nContentType: "+schemas.InvalidContentTypeError+"\nOwnerID: "+schemas.InvalidFileOwnerIDError+"\nMetadata: "+schemas.InvalidFileMetadataError, "Error message should match")

		file = &UserFile{
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

		err = file.ShouldHaveFields("john_doe")(schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Directory, schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.Metadata)
		assert.NoError(t, err, "ShouldHaveFields should not return an error for valid fields")

		err = file.ShouldHaveFields("john_doe")()
		assert.Error(t, err, "ShouldHaveFields should return an error when no fields are specified")
		assert.EqualError(t, err, schemas.NoFieldsSpecifiedError, "Error message should match")

		err = file.ShouldHaveFields("john_doe")("invalid_field")
		assert.Error(t, err, "ShouldHaveFields should return an error for invalid field")
		assert.EqualError(t, err, schemas.ValidationErrorPrefix+"field 'invalid_field' is not a valid file field", "Error message should match")
	})
}
