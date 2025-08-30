package foltests

import (
	"context"
	"errors"
	"testing"
	"time"

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
	fvp := folschema.New().
		WithName("newdocs").
		WithPath("john_doe/newdocs")

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

	fol, err := repo.UpdateFolder(context.Background(), folder, fvp, owner)
	assert.NoError(t, err)
	assert.Equal(t, fvp.Name, fol.Name)
	assert.Equal(t, fvp.Path, fol.Path)
	assert.Equal(t, 2, *fol.Parent)
}

func TestUpdateFolder_RejectIDUpdate(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)
	folder := makeUserFolder()
	owner := makeOwner()
	fvp := folschema.New().
		WithID(folder.ID)

	f, err := repo.UpdateFolder(context.Background(), folder, fvp, owner)
	assert.Error(t, err)
	assert.Nil(t, f, "folder should be nil")
	assert.Contains(t, err.Error(), "id should not be set")
}

func TestUpdateFolder_RejectOwnerIDUpdate(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)
	folder := makeUserFolder()
	owner := makeOwner()
	fvp := folschema.New().
		WithOwnerID(owner.ID)

	f, err := repo.UpdateFolder(context.Background(), folder, fvp, owner)
	assert.Error(t, err)
	assert.Nil(t, f, "folder should be nil")
	assert.Contains(t, err.Error(), "owner_id should not be set")
}

func TestUpdateFolder_RejectCreatedAtUpdate(t *testing.T) {
	mockDB := &db.Database{}
	repo := folderrepo.NewFolderRepo(mockDB)
	folder := makeUserFolder()
	owner := makeOwner()
	fvp := &folschema.UserFolder{
		CreatedAt: time.Now(),
	}

	f, err := repo.UpdateFolder(context.Background(), folder, fvp, owner)
	assert.Error(t, err)
	assert.Nil(t, f, "folder should be nil")
	assert.Contains(t, err.Error(), "created_at should not be set")
}

func TestUpdateFolder_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	folder := makeUserFolder()
	owner := makeOwner()
	fvp := folschema.New().
		WithName("name")

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("error updating folder"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	f, err := repo.UpdateFolder(context.Background(), folder, fvp, owner)
	assert.Error(t, err)
	assert.Nil(t, f, "folder should be nil")
	assert.Contains(t, err.Error(), "error updating folder")
}
