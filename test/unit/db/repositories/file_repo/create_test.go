package filerepo_test

import (
	"context"
	"errors"
	"testing"

	"homeserver/db"
	filerepo "homeserver/db/repositories/file_repo"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/test/helpers/mocks"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fieldDescs = []pgconn.FieldDescription{
	{Name: "id"},
	{Name: "filename"},
	{Name: "extension"},
	{Name: "path"},
	{Name: "directory_id"},
	{Name: "size"},
	{Name: "content_type"},
	{Name: "owner_id"},
	{Name: "created_at"},
	{Name: "updated_at"},
	{Name: "metadata"},
}

func TestCreateFile_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	expected := makeUserFile()

	owner := makeOwner()
	file := &fschema.UserFile{
		Filename:      "file",
		Extension:     "txt",
		Path:          "john_doe/docs",
		Directory:     1,
		Size:          1234,
		ContentType:   "text/plain",
		TrustedSource: []byte("trusted"),
		OwnerID:       42,
	}

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, []any{expected.ID, expected.Filename, expected.Extension, expected.Path, expected.Directory, expected.Size, expected.ContentType, expected.OwnerID})
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	created, err := repo.CreateFile(context.Background(), file, owner)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, expected.ID, created.ID)
	assert.Equal(t, expected.Filename, created.Filename)
	assert.Equal(t, expected.Extension, created.Extension)
	assert.Equal(t, expected.Path, created.Path)
	assert.Equal(t, expected.Size, created.Size)
	assert.Equal(t, expected.ContentType, created.ContentType)
	assert.Equal(t, expected.OwnerID, created.OwnerID)
}

func TestCreateFile_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(errors.New("DB error"))

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	owner := makeOwner()
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
