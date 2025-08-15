package filerepo_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"homeserver/db"
	filerepo "homeserver/db/repositories/file_repo"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/test"
	"homeserver/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFuncs struct {
	mock.Mock
}

func (m *MockFuncs) FixDuplicates(f *fschema.UserFile) (*fschema.UserFile, error) {
	args := m.Called(f)
	return args.Get(0).(*fschema.UserFile), args.Error(1)
}

func TestUpdateFileEntry_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := makeUserFile()
	owner := makeOwner()
	fieldsToUpdate := []string{"filename", "extension", "path"}

	mockFuncs := &MockFuncs{}
	mockFuncs.On("FixDuplicates", mock.Anything).Return(file, nil)
	release := test.MockFunction(&filerepo.FixDuplicates, mockFuncs.FixDuplicates, func() {})
	defer release()

	newValues := []any{
		file.ID + 1, file.Filename, file.Extension, file.Path, file.Directory,
		file.Size, file.ContentType, file.OwnerID, file.CreatedAt, file.UpdatedAt,
	}

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, newValues)
	})

	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	err := repo.UpdateFileEntry(context.Background(), file, fieldsToUpdate, owner)
	assert.NoError(t, err)
	assert.NotEqual(t, file.ID, 1)
}

func TestUpdateFileEntry_PreventDuplicateError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := makeUserFile()
	owner := makeOwner()
	fieldsToUpdate := []string{"filename", "extension", "path"}

	// Patch FixDuplicates to return error
	mockFuncs := &MockFuncs{}
	mockFuncs.On("FixDuplicates", mock.Anything).Return((*fschema.UserFile)(nil), errors.New("duplicate error"))
	release := test.MockFunction(&filerepo.FixDuplicates, mockFuncs.FixDuplicates, func() {})
	defer release()

	err := repo.UpdateFileEntry(context.Background(), file, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preventing duplicate file name")
}

func TestUpdateFileEntry_PartialNameFields(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := &fschema.UserFile{
		ID:       1,
		Filename: "file",
	}
	owner := makeOwner()
	fieldsToUpdate := []string{"filename"}

	err := repo.UpdateFileEntry(context.Background(), file, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "if filename, extension, or path is provided for file update, all must be provided")
}

func TestUpdateFileEntry_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := makeUserFile()
	owner := makeOwner()
	fieldsToUpdate := []string{"filename", "extension", "path"}

	mockFuncs := &MockFuncs{}
	mockFuncs.On("FixDuplicates", mock.Anything).Return(file, nil)
	release := test.MockFunction(&filerepo.FixDuplicates, mockFuncs.FixDuplicates, func() {})
	defer release()

	mockRow := &mocks.MockRow{}
	mockRow.On("Scan", mock.Anything).Return(fmt.Errorf("update error"))
	mockPool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRow)

	err := repo.UpdateFileEntry(context.Background(), file, fieldsToUpdate, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating file entry")
}
