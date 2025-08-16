package foltests

import (
	"context"
	"errors"
	"testing"

	"homeserver/db"
	folderrepo "homeserver/db/repositories/folder_repo"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetFolderByID_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	folder, err := repo.GetFolderByID(context.Background(), expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, folder.ID)
	assert.Equal(t, expected.Name, folder.Name)
	assert.Equal(t, expected.Path, folder.Path)
	assert.Equal(t, *expected.Parent, *folder.Parent)
	assert.Equal(t, expected.OwnerID, folder.OwnerID)
}

func TestGetFolderByID_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	folder, err := repo.GetFolderByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, folder)
}

func TestGetFolderByPath_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	folder, err := repo.GetFolderByPath(context.Background(), expected)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, folder.ID)
	assert.Equal(t, expected.Name, folder.Name)
	assert.Equal(t, expected.Path, folder.Path)
	assert.Equal(t, *expected.Parent, *folder.Parent)
	assert.Equal(t, expected.OwnerID, folder.OwnerID)
}

func TestGetFolderByPath_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	folder, err := repo.GetFolderByPath(context.Background(), &folschema.UserFolder{Path: "p", Name: "n"})
	assert.Error(t, err)
	assert.Nil(t, folder)
}

func TestGetFolderByOwner_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	folder, err := repo.GetFolderByOwner(context.Background(), expected.OwnerID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, folder.ID)
	assert.Equal(t, expected.Name, folder.Name)
	assert.Equal(t, expected.Path, folder.Path)
	assert.Equal(t, *expected.Parent, *folder.Parent)
	assert.Equal(t, expected.OwnerID, folder.OwnerID)
}

func TestGetFolderByOwner_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	folder, err := repo.GetFolderByOwner(context.Background(), 42)
	assert.Error(t, err)
	assert.Nil(t, folder)
}

func TestGetChildren_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRows := &mocks.MockRows{}
	// Simulate one row, then end
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil)
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})
	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	children, err := repo.GetChildren(context.Background(), *expected.Parent)
	assert.NoError(t, err)
	assert.Len(t, children, 1)
	assert.Equal(t, expected.ID, children[0].ID)
	assert.Equal(t, expected.Name, children[0].Name)
	assert.Equal(t, expected.Path, children[0].Path)
	assert.Equal(t, *expected.Parent, *children[0].Parent)
	assert.Equal(t, expected.OwnerID, children[0].OwnerID)
}

func TestGetChildren_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(&mocks.MockRows{}, errors.New("db error"))

	children, err := repo.GetChildren(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, children)
}
