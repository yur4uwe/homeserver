package foltests

import (
	"context"
	"errors"
	"testing"
	"time"

	"homeserver/db"
	folderrepo "homeserver/db/repositories/folder_repo"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateFolder_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()

	expected := &folschema.UserFolder{
		ID:        10,
		Name:      folder.Name,
		Path:      folder.Path,
		Parent:    folder.Parent,
		OwnerID:   folder.OwnerID,
		CreatedAt: time.Now(),
	}

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID, expected.CreatedAt}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	inserted, err := repo.CreateFolder(context.Background(), folder, owner)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, inserted.ID)
	assert.Equal(t, expected.Name, inserted.Name)
	assert.Equal(t, expected.Path, inserted.Path)
	assert.Equal(t, *expected.Parent, *inserted.Parent)
	assert.Equal(t, expected.OwnerID, inserted.OwnerID)
}

func TestCreateFolder_ValidationError(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()
	folder.Name = "" // Should trigger validation error

	inserted, err := repo.CreateFolder(context.Background(), folder, owner)
	assert.Error(t, err)
	assert.Nil(t, inserted)
}

func TestCreateFolder_InvalidParent(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()
	invalidParent := 0
	folder.Parent = &invalidParent

	inserted, err := repo.CreateFolder(context.Background(), folder, owner)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "invalid parent folder ID")
}

func TestCreateFolder_PathParentMismatch(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)

	owner := makeOwner()

	// Path set, parent nil
	folder1 := makeUserFolder()
	folder1.Parent = nil
	folder1.Path = "some/path"
	inserted, err := repo.CreateFolder(context.Background(), folder1, owner)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "folder path and parent folder must be set together")

	// Path empty, parent set
	folder2 := makeUserFolder()
	folder2.Path = ""
	inserted, err = repo.CreateFolder(context.Background(), folder2, owner)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "folder path and parent folder must be set together")
}

func TestCreateFolder_RootFolderNameMismatch(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)

	owner := makeOwner()
	folder := &folschema.UserFolder{
		Name:    "not_owner_name",
		Path:    "",
		Parent:  nil,
		OwnerID: owner.ID,
	}

	inserted, err := repo.CreateFolder(context.Background(), folder, owner)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "root folder must have the same name as the owner")
}

func TestCreateFolder_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	inserted, err := repo.CreateFolder(context.Background(), folder, owner)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Contains(t, err.Error(), "db error")
}
