package filerepo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"homeserver/db"
	filerepo "homeserver/db/repositories/file_repo"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateFile_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockRow := mocks.MockRowReturn(
		1,
		"file",
		"txt",
		"john_doe/docs",
		1,
		int64(1234),
		"text/plain",
		42,
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		"some metadata",
	)

	mockPool.On("QueryRow",
		mock.AnythingOfType("backgroundCtx"),
		mock.AnythingOfType("string"),
		mock.Anything).Return(mockRow)

	owner := &uschema.User{ID: 42, Username: "john_doe"}
	file := &fschema.UserFile{
		Filename:    "file",
		Extension:   "txt",
		Path:        "john_doe/docs",
		Directory:   1,
		Size:        1234,
		ContentType: "text/plain",
		OwnerID:     42,
	}

	created, err := repo.CreateFile(context.Background(), file, owner)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, "file", created.Filename)
	assert.Equal(t, "txt", created.Extension)
	assert.Equal(t, "john_doe/docs", created.Path)
	assert.Equal(t, int64(1234), created.Size)
	assert.Equal(t, "text/plain", created.ContentType)
	assert.Equal(t, 42, created.OwnerID)
}

func TestCreateFile_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(fmt.Errorf("DB error"))

	mockPool.On("QueryRow",
		mock.AnythingOfType("backgroundCtx"),
		mock.AnythingOfType("string"),
		mock.Anything).Return(mockRow)

	owner := &uschema.User{ID: 42, Username: "john_doe"}
	file := &fschema.UserFile{
		Filename:    "file.txt",
		Extension:   "txt",
		Path:        "john_doe/docs",
		Directory:   1,
		Size:        1234,
		ContentType: "text/plain",
		OwnerID:     42,
	}

	created, err := repo.CreateFile(context.Background(), file, owner)
	assert.Error(t, err)
	assert.Nil(t, created)
}
