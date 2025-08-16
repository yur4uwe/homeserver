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

func TestDeleteFolderByID_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expectedID := 1
	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = expectedID
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	deletedID, err := repo.DeleteFolderByID(context.Background(), expectedID, 42)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, deletedID)
}

func TestDeleteFolderByID_InvalidID(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)

	deletedID, err := repo.DeleteFolderByID(context.Background(), 0, 42)
	assert.Error(t, err)
	assert.Equal(t, 0, deletedID)
	assert.Contains(t, err.Error(), "invalid folder ID")
}

func TestDeleteFolderByID_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	deletedID, err := repo.DeleteFolderByID(context.Background(), 1, 42)
	assert.Error(t, err)
	assert.Equal(t, 0, deletedID)
	assert.Contains(t, err.Error(), "error deleting folder entry")
}

func TestDeleteFolderByPath_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()
	expectedID := 1

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		*dest[0].(*int) = expectedID
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	deletedID, err := repo.DeleteFolderByPath(context.Background(), folder, owner)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, deletedID)
}

func TestDeleteFolderByPath_InvalidFolderData(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := &folschema.UserFolder{Path: "", Name: "", OwnerID: 0}
	owner := makeOwner()

	deletedID, err := repo.DeleteFolderByPath(context.Background(), folder, owner)
	assert.Error(t, err)
	assert.Equal(t, 0, deletedID)
	assert.Contains(t, err.Error(), "invalid folder data")
}

func TestDeleteFolderByPath_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	deletedID, err := repo.DeleteFolderByPath(context.Background(), folder, owner)
	assert.Error(t, err)
	assert.Equal(t, 0, deletedID)
	assert.Contains(t, err.Error(), "error deleting folder entry")
}
