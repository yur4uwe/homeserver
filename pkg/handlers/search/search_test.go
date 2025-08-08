package search

import (
	"context"
	"errors"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/test"
	"homeserver/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSearchHandler(t *testing.T) {
	mockDB := new(mocks.MockDB)
	requester := &user_schema.User{ID: 1, Username: "testuser"}
	ctx := context.Background()
	query := "report"
	directory := []string{"testuser", "docs"}

	file := &file_schema.UserFile{
		ID:        1,
		Filename:  "report",
		Extension: "pdf",
		Path:      "testuser/docs",
		Size:      2048,
	}

	resetter := test.NewMocksResetter(&mockDB.Mock)

	t.Run("success", func(t *testing.T) {
		resetter()

		mockDB.On("SearchFiles", ctx, query, requester, directory).Return([]*file_schema.UserFile{file}, nil)
		data, err := FileSearchHandler(ctx, mockDB, query, directory, requester)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.True(t, data.QueryResults)
		assert.Len(t, data.Entries, 1)
		assert.Equal(t, "report", data.Entries[0].Name)
	})

	t.Run("error searching files", func(t *testing.T) {
		resetter()

		mockDB.On("SearchFiles", ctx, query, requester, directory).Return(([]*file_schema.UserFile)(nil), errors.New("fail search"))
		data, err := FileSearchHandler(ctx, mockDB, query, directory, requester)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "error searching files")
	})
}
