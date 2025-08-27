package search_test

import (
	"context"
	"errors"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/pkg/handlers/web/search"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSearchHandler(t *testing.T) {
	_, mockFR, _, _, _ := mocks.InitRepoMocks()
	mockDB := repositories.NewRepositories(nil, mockFR, nil, nil, nil)
	requester := &uschema.User{ID: 1, Username: "testuser"}
	ctx := context.Background()
	query := "report"
	directory := []string{"testuser", "docs"}

	file := fschema.UserFile{
		ID:        1,
		Filename:  "report",
		Extension: "pdf",
		Path:      "testuser/docs",
		Size:      2048,
	}

	resetter := helpers.NewMocksResetter(&mockFR.Mock)

	t.Run("success", func(t *testing.T) {
		resetter()

		mockFR.On("FuzzySearch", ctx, query, strings.Join(directory, "/"), requester).Return([]fschema.UserFile{file}, nil)
		data, err := search.FileSearchHandler(ctx, mockDB, query, directory, requester)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.True(t, data.QueryResults)
		assert.Len(t, data.Entries, 1)
		assert.Equal(t, "report", data.Entries[0].Name)
	})

	t.Run("error searching files", func(t *testing.T) {
		resetter()

		mockFR.On("FuzzySearch", ctx, query, strings.Join(directory, "/"), requester).Return(([]fschema.UserFile)(nil), errors.New("fail search"))
		data, err := search.FileSearchHandler(ctx, mockDB, query, directory, requester)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "error searching files")
	})
}
