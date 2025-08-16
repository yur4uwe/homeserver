package filerepo_test

import (
	"homeserver/constants"
	"homeserver/constants/fconsts"
	filerepo "homeserver/db/repositories/file_repo"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/test"
	"homeserver/test/mocks"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreventDuplicateFileNames(t *testing.T) {
	mockOS := new(mocks.MockOS)

	restore := test.MockFunction(&filerepo.Stat, mockOS.Stat, func() {})
	defer restore()
	restoreGetenv := test.MockFunction(&fschema.Getenv, mockOS.Getenv, func() {})
	defer restoreGetenv()

	resetter := test.NewMocksResetter(
		&mockOS.Mock,
	)

	t.Run("should return the file if it does not exist", func(t *testing.T) {
		resetter()
		mockOS.On("Getenv", constants.StoragePath).Return("/mock/storage/path/")
		mockOS.On("Stat", "/mock/storage/path/john_doe/docs/file.txt").Return(nil, os.ErrNotExist)

		mockFile := fschema.NewUserFileBuilder().
			WithFullName("file.txt").
			WithPath("john_doe/docs")

		file, err := filerepo.PreventDuplicateFileName(mockFile)
		assert.NoError(t, err)
		assert.False(t, file == mockFile, "Returned file should not be the same instance as the original")
		assert.Equal(t, mockFile.Filename, file.Filename)
	})

	t.Run("should append a counter to the filename if it exists", func(t *testing.T) {
		resetter()
		mockOS.On("Getenv", constants.StoragePath).Return("/mock/storage/path/")
		mockOS.On("Stat", "/mock/storage/path/john_doe/docs/file.txt").Return(nil, nil)
		mockOS.On("Stat", "/mock/storage/path/john_doe/docs/file_1.txt").Return(nil, os.ErrNotExist)

		mockFile := fschema.NewUserFileBuilder().
			WithFullName("file.txt").
			WithPath("john_doe/docs")

		file, err := filerepo.PreventDuplicateFileName(mockFile)
		assert.NoError(t, err)
		assert.False(t, file == mockFile, "Returned file should not be the same instance as the original")
		assert.NotEqual(t, mockFile.Filename, file.Filename)
		assert.Equal(t, "file_1", file.Filename)
	})
}

func TestOrderFileUpdateFields(t *testing.T) {
	metadata := "meta"
	file := &fschema.UserFile{
		Filename:    "file.txt",
		Extension:   "txt",
		Path:        "john_doe/docs",
		Directory:   1,
		Size:        1234,
		ContentType: "text/plain",
		OwnerID:     42,
		Metadata:    &metadata,
	}

	t.Run("all fields", func(t *testing.T) {
		fields := []string{
			fconsts.Filename, fconsts.Extension, fconsts.Path, fconsts.Directory,
			fconsts.Size, fconsts.ContentType, fconsts.OwnerID, fconsts.Metadata,
		}
		args := filerepo.OrderFileUpdateFields(fields, file)
		expected := []any{
			file.Filename, file.Extension, file.Path, file.Directory,
			file.Size, file.ContentType, file.OwnerID, file.Metadata,
		}
		assert.Equal(t, expected, args)
	})

	t.Run("single field", func(t *testing.T) {
		fields := []string{fconsts.Filename}
		args := filerepo.OrderFileUpdateFields(fields, file)
		expected := []any{file.Filename}
		assert.Equal(t, expected, args)
	})

	t.Run("mixed order", func(t *testing.T) {
		fields := []string{fconsts.Size, fconsts.Filename, fconsts.Directory}
		args := filerepo.OrderFileUpdateFields(fields, file)
		expected := []any{file.Size, file.Filename, file.Directory}
		assert.Equal(t, expected, args)
	})

	t.Run("empty fields", func(t *testing.T) {
		fields := []string{}
		args := filerepo.OrderFileUpdateFields(fields, file)
		assert.Empty(t, args)
	})
}
