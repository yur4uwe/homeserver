package filerepo_test

import (
	"context"
	"errors"
	"testing"

	"homeserver/db"
	filerepo "homeserver/db/repositories/file_repo"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteFileByID_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)
	owner := &uschema.User{ID: 42}

	tag := pgconn.NewCommandTag("DELETE 1")

	mockPool.On("Exec",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(tag, nil)

	err := repo.DeleteFileByID(context.Background(), 123, owner)
	assert.NoError(t, err)
}

func TestDeleteFileByID_NotFound(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)
	owner := &uschema.User{ID: 42}

	tag := pgconn.NewCommandTag("DELETE 0")

	mockPool.On("Exec",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(tag, nil)

	err := repo.DeleteFileByID(context.Background(), 123, owner)
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to delete file")
}

func TestDeleteFileByID_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)
	owner := &uschema.User{ID: 42}

	tag := pgconn.NewCommandTag("DELETE 0")

	mockPool.On("Exec",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(tag, errors.New("db error"))

	err := repo.DeleteFileByID(context.Background(), 123, owner)
	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
}

func TestDeleteFile_ValidationError(t *testing.T) {
	mockDB := &db.Database{}
	repo := filerepo.NewFileRepo(mockDB)
	owner := &uschema.User{ID: 0} // Invalid ID

	file := &fschema.UserFile{Filename: "file", Path: "path"}
	_, err := repo.DeleteFile(context.Background(), file, owner)
	assert.Error(t, err)
}

func TestDeleteFile_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)
	owner := &uschema.User{ID: 42, Username: "john_doe"}
	file := &fschema.UserFile{Filename: "file", Path: "john_doe/docs"}

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = 1
	})

	mockPool.On("QueryRow",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mockRow)

	id, err := repo.DeleteFile(context.Background(), file, owner)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestDeleteFile_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)
	owner := &uschema.User{ID: 42, Username: "john_doe"}
	file := &fschema.UserFile{Filename: "file", Path: "john_doe/docs"}

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))

	mockPool.On("QueryRow",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mockRow)

	id, err := repo.DeleteFile(context.Background(), file, owner)
	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.EqualError(t, err, "db error")
}
