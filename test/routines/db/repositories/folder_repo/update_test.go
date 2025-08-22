package foltests

import (
	"context"
	"errors"
	"testing"
	"time"

	"homeserver/constants/folconsts"
	"homeserver/db"
	folderrepo "homeserver/db/repositories/folder_repo"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/helpers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeUserFolder() *folschema.UserFolder {
	parent := 1
	return &folschema.UserFolder{
		ID:        1,
		Name:      "docs",
		Path:      "john_doe/docs",
		Parent:    &parent,
		OwnerID:   42,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func makeOwner() *uschema.User {
	return &uschema.User{ID: 42, Username: "john_doe"}
}

func TestUpdateFolder_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()
	fieldsToUpdate := []string{folconsts.FolderName, folconsts.Path}

	newParent := 2
	newValues := []any{
		folder.ID, "newdocs", "john_doe/newdocs", &newParent, folder.OwnerID, folder.CreatedAt,
	}

	mockRows := new(mocks.MockRows)
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, newValues)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	err := repo.UpdateFolder(context.Background(), folder, fieldsToUpdate, owner)
	assert.NoError(t, err)
	assert.Equal(t, "newdocs", folder.Name)
	assert.Equal(t, "john_doe/newdocs", folder.Path)
	assert.Equal(t, 2, *folder.Parent)
}

func TestUpdateFolder_RejectIDUpdate(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)
	folder := makeUserFolder()
	owner := makeOwner()
	fieldsToUpdate := []string{folconsts.ID}

	err := repo.UpdateFolder(context.Background(), folder, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot update folder ID")
}

func TestUpdateFolder_RejectOwnerIDUpdate(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)
	folder := makeUserFolder()
	owner := makeOwner()
	fieldsToUpdate := []string{folconsts.OwnerID}

	err := repo.UpdateFolder(context.Background(), folder, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot update owner ID")
}

func TestUpdateFolder_RejectCreatedAtUpdate(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)
	folder := makeUserFolder()
	owner := makeOwner()
	fieldsToUpdate := []string{folconsts.CreatedAt}

	err := repo.UpdateFolder(context.Background(), folder, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot update created at")
}

func TestUpdateFolder_InvalidFolderData(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)
	folder := makeUserFolder()
	owner := makeOwner()
	fieldsToUpdate := []string{folconsts.FolderName, folconsts.Path}
	folder.Name = ""

	err := repo.UpdateFolder(context.Background(), folder, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid folder data")
}

func TestUpdateFolder_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()
	fieldsToUpdate := []string{folconsts.FolderName, folconsts.Path}

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("error updating folder"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	err := repo.UpdateFolder(context.Background(), folder, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating folder")
}
