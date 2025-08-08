package search

import (
	"context"
	"errors"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/test"
	"homeserver/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrowseHandler(t *testing.T) {
	mockDB := new(mocks.MockDB)
	requester := &user_schema.User{ID: 1, Username: "testuser"}
	ctx := context.Background()
	directory := []string{"testuser", "docs"}

	resetter := test.NewMocksResetter(&mockDB.Mock)

	file := file_schema.NewUserFileBuilder().
		WithID(1).
		WithFilename("file1").
		WithExtension("txt").
		WithPath("testuser/docs/file1.txt").
		WithSize(100)

	folder := folder_schema.NewUserFolderBuilder().
		WithID(2).
		WithName("folder1").
		WithPath("testuser/docs/folder1")

	t.Run("success", func(t *testing.T) {
		resetter()

		mockDB.On("ExtractFilesFromDirectory", ctx, "testuser/docs", requester).Return([]*file_schema.UserFile{file}, nil)
		mockDB.On("ExtractFoldersFromDirectory", ctx, directory, requester).Return([]*folder_schema.UserFolder{folder}, nil)

		data, err := BrowseHandler(ctx, mockDB, directory, requester)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, "docs", data.Directory)
		assert.Equal(t, "docs", data.URIDirectory)
		assert.False(t, data.QueryResults)
		assert.Len(t, data.Entries, 2)
		// Folders should come first
		assert.True(t, data.Entries[0].IsFolder)
		assert.Equal(t, "folder1", data.Entries[0].Name)
		assert.False(t, data.Entries[1].IsFolder)
		assert.Equal(t, "file1", data.Entries[1].Name)
	})

	t.Run("error extracting files", func(t *testing.T) {
		resetter()

		mockDB.On("ExtractFilesFromDirectory", ctx, "testuser/docs", requester).Return(([]*file_schema.UserFile)(nil), errors.New("fail files"))
		mockDB.On("ExtractFoldersFromDirectory", ctx, directory, requester).Return(([]*folder_schema.UserFolder)(nil), nil)

		data, err := BrowseHandler(ctx, mockDB, directory, requester)

		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "error extracting user files")
	})

	t.Run("error extracting folders", func(t *testing.T) {
		resetter()

		mockDB.On("ExtractFilesFromDirectory", ctx, "testuser/docs", requester).Return([]*file_schema.UserFile{}, nil)
		mockDB.On("ExtractFoldersFromDirectory", ctx, directory, requester).Return(([]*folder_schema.UserFolder)(nil), errors.New("fail folders"))

		data, err := BrowseHandler(ctx, mockDB, directory, requester)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "error extracting user folders")
	})
}
