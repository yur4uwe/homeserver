package foltests

import (
	"context"
	"errors"
	"strings"
	"testing"

	"homeserver/db"
	folderrepo "homeserver/db/repositories/folder_repo"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/test/helpers/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fieldDescs = []pgconn.FieldDescription{
	{Name: "id"},
	{Name: "name"},
	{Name: "path"},
	{Name: "parent"},
	{Name: "owner_id"},
	{Name: "created_at"},
}

func TestGetFolder_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		vals := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID, expected.CreatedAt}
		mocks.RunScan(dest, vals)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	folder, err := repo.GetFolder(context.Background(), folschema.New().WithID(expected.ID))
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

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("db error"))
	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	folder, err := repo.GetFolder(context.Background(), folschema.New().WithID(1))
	assert.Error(t, err)
	assert.Nil(t, folder)
}

func TestGetFolderByPath_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		vals := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID, expected.CreatedAt}
		mocks.RunScan(dest, vals)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	folder, err := repo.GetFolder(context.Background(), expected)
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

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("db error"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	folder, err := repo.GetFolder(context.Background(), folschema.New().WithPath("p").WithName("n"))
	assert.Error(t, err)
	assert.Nil(t, folder)
}

func TestGetFolderByOwner_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		vals := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID, expected.CreatedAt}
		mocks.RunScan(dest, vals)
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	folder, err := repo.GetFoldersByOwner(context.Background(), expected.OwnerID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, folder[0].ID)
	assert.Equal(t, expected.Name, folder[0].Name)
	assert.Equal(t, expected.Path, folder[0].Path)
	assert.Equal(t, *expected.Parent, *folder[0].Parent)
	assert.Equal(t, expected.OwnerID, folder[0].OwnerID)
}

func TestGetFolderByOwner_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("db error"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	folder, err := repo.GetFoldersByOwner(context.Background(), 42)
	assert.Error(t, err)
	assert.Nil(t, folder)
}

func TestGetChildren_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := folderrepo.NewFolderRepo(mockDB)

	expected := makeUserFolder()
	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil)
	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		values := []any{expected.ID, expected.Name, expected.Path, expected.Parent, expected.OwnerID}
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, values)
	})
	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	mockRequester := &uschema.User{ID: expected.OwnerID, Username: "john_doe"}

	pathSlice := strings.Split(expected.Path, "/")
	name := pathSlice[len(pathSlice)-1]
	path := strings.Join(pathSlice[:len(pathSlice)-1], "/")

	children, err := repo.GetChildren(context.Background(), path, name, mockRequester)
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

	pathSlice := strings.Split("john_doe/folder", "/")
	name := pathSlice[len(pathSlice)-1]
	path := strings.Join(pathSlice[:len(pathSlice)-1], "/")

	mockRequester := &uschema.User{ID: 1, Username: "john_doe"}

	children, err := repo.GetChildren(context.Background(), path, name, mockRequester)
	assert.Error(t, err)
	assert.Nil(t, children)
}
