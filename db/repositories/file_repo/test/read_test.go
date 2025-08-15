package filerepo_test

import (
	"context"
	"errors"
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

func makeUserFile() *fschema.UserFile {
	metadata := "meta"
	return &fschema.UserFile{
		ID:          1,
		Filename:    "file",
		Extension:   "txt",
		Path:        "john_doe/docs",
		Directory:   1,
		Size:        1234,
		ContentType: "text/plain",
		OwnerID:     42,
		CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Metadata:    &metadata,
	}
}

func makeOwner() *uschema.User {
	return &uschema.User{ID: 42, Username: "john_doe"}
}

func TestGetFileByID_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	expected := makeUserFile()
	mockRow := mocks.MockRowReturn(
		expected.ID, expected.Filename, expected.Extension, expected.Path, expected.Directory,
		expected.Size, expected.ContentType, expected.OwnerID, expected.CreatedAt, expected.UpdatedAt, *expected.Metadata,
	)
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	file, err := repo.GetFileByID(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, expected.Filename, file.Filename)
	assert.Equal(t, expected.Path, file.Path)
}

func TestGetFileByID_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))

	file, err := repo.GetFileByID(context.Background(), "1")
	assert.Error(t, err)
	assert.Nil(t, file)
}

func TestGetFullFileInfo_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	expected := makeUserFile()
	mockRow := mocks.MockRowReturn(
		expected.ID, expected.Filename, expected.Extension, expected.Path, expected.Directory,
		expected.Size, expected.ContentType, expected.OwnerID, expected.CreatedAt, expected.UpdatedAt, *expected.Metadata,
	)
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	owner := makeOwner()
	file, err := repo.GetFullFileInfo(context.Background(), expected, owner)
	assert.NoError(t, err)
	assert.Equal(t, expected.Filename, file.Filename)
	assert.Equal(t, expected.Path, file.Path)
}

func TestGetFullFileInfo_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))

	owner := makeOwner()
	file := makeUserFile()
	result, err := repo.GetFullFileInfo(context.Background(), file, owner)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetFilesByUser_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	expected := makeUserFile()
	expected_values := []any{
		expected.ID, expected.Filename, expected.Extension, expected.Path, expected.Directory,
		expected.Size, expected.ContentType, expected.OwnerID, expected.CreatedAt, expected.UpdatedAt, *expected.Metadata,
	}
	mockRows := &mocks.MockRows{}
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil)
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		destSlice := args.Get(0).([]any)
		mocks.RunScan(destSlice, expected_values)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	owner := makeOwner()
	files, err := repo.GetFilesByUser(context.Background(), owner)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, expected.Filename, files[0].Filename)
}

func TestGetFilesByUser_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(&mocks.MockRows{}, errors.New("db error"))

	owner := makeOwner()
	files, err := repo.GetFilesByUser(context.Background(), owner)
	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestFuzzySearch_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	expected := makeUserFile()
	expected_values := []any{
		expected.ID, expected.Filename, expected.Extension, expected.Path, expected.Directory,
		expected.Size, expected.ContentType, expected.OwnerID, expected.CreatedAt, expected.UpdatedAt, *expected.Metadata,
	}
	mockRows := &mocks.MockRows{}
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil)
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		destSlice := args.Get(0).([]any)
		mocks.RunScan(destSlice, expected_values)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(mockRows, nil)

	owner := makeOwner()
	files, err := repo.FuzzySearch(context.Background(), "file", "john_doe/docs", owner)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, expected.Filename, files[0].Filename)
}

func TestFuzzySearch_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(&mocks.MockRows{}, errors.New("db error"))

	owner := makeOwner()
	files, err := repo.FuzzySearch(context.Background(), "file", "john_doe/docs", owner)
	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestGetFilesByDirectory_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	expected := makeUserFile()
	expected_values := []any{
		expected.ID, expected.Filename, expected.Extension, expected.Path, expected.Directory,
		expected.Size, expected.ContentType, expected.OwnerID, expected.CreatedAt, expected.UpdatedAt, *expected.Metadata,
	}
	mockRows := &mocks.MockRows{}
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil)
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		destSlice := args.Get(0).([]any)
		mocks.RunScan(destSlice, expected_values)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(mockRows, nil)

	owner := makeOwner()
	files, err := repo.GetFilesByDirectory(context.Background(), "john_doe/docs", owner)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, expected.Filename, files[0].Filename)
}

func TestGetFilesByDirectory_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(&mocks.MockRows{}, errors.New("db error"))

	owner := makeOwner()
	files, err := repo.GetFilesByDirectory(context.Background(), "john_doe/docs", owner)
	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestGetFiles_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	expected := makeUserFile()
	expected_values := []any{
		expected.ID, expected.Filename, expected.Extension, expected.Path, expected.Directory,
		expected.Size, expected.ContentType, expected.OwnerID, expected.CreatedAt, expected.UpdatedAt, *expected.Metadata,
	}
	mockRows := &mocks.MockRows{}
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil)
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		destSlice := args.Get(0).([]any)
		mocks.RunScan(destSlice, expected_values)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	files, err := repo.GetFiles(context.Background())
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, expected.Filename, files[0].Filename)
}

func TestGetFiles_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(&mocks.MockRows{}, errors.New("db error"))

	files, err := repo.GetFiles(context.Background())
	assert.Error(t, err)
	assert.Nil(t, files)
}
