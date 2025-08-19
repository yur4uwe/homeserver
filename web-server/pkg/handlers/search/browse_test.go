package search

import (
	"context"
	"errors"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test"
	"homeserver/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrowseHandler(t *testing.T) {
	_, mockFR, mockFolR, _, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, mockFR, mockFolR, nil, nil)
	requester := &uschema.User{ID: 1, Username: "testuser"}
	ctx := context.Background()
	directory := []string{"testuser", "docs"}

	resetter := test.NewMocksResetter(&mockFR.Mock, &mockFolR.Mock)

	file := fschema.NewUserFileBuilder().
		WithID(1).
		WithFilename("file1").
		WithExtension("txt").
		WithPath("testuser/docs/file1.txt").
		WithSize(100)

	folder := folschema.NewUserFolderBuilder().
		WithID(2).
		WithName("folder1").
		WithPath("testuser/docs/folder1")

	t.Run("success", func(t *testing.T) {
		resetter()

		mockFR.On("GetFilesByDirectory", ctx, "testuser/docs", requester).Return([]fschema.UserFile{*file}, nil)
		mockFolR.On("GetChildren", ctx, "testuser", "docs", requester).Return([]folschema.UserFolder{*folder}, nil)

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

		mockFR.On("GetFilesByDirectory", ctx, "testuser/docs", requester).Return(([]fschema.UserFile)(nil), errors.New("fail files"))
		mockFolR.On("GetChildren", ctx, "testuser", "docs", requester).Return(([]folschema.UserFolder)(nil), nil)

		data, err := BrowseHandler(ctx, mockDB, directory, requester)

		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "error extracting user files")
	})

	t.Run("error extracting folders", func(t *testing.T) {
		resetter()

		mockFR.On("GetFilesByDirectory", ctx, "testuser/docs", requester).Return(([]fschema.UserFile)(nil), nil)
		mockFolR.On("GetChildren", ctx, "testuser", "docs", requester).Return(([]folschema.UserFolder)(nil), errors.New("fail folders"))

		data, err := BrowseHandler(ctx, mockDB, directory, requester)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "error extracting user folders")
	})
}
