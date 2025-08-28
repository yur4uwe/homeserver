package filerepo_test

import (
	"context"
	"fmt"
	"testing"

	"homeserver/db"
	filerepo "homeserver/db/repositories/file_repo"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks"

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

var mockGetFileByID = func(repo *filerepo.FileRepo, ctx context.Context, id int) (*fschema.UserFile, error) {
	// Return a dummy file or whatever you need for the test
	return makeUserFile(), nil
}

func TestUpdateFileEntry_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := makeUserFile()
	owner := makeOwner()
	fvp := fschema.NewUserFileBuilder().
		WithFilename("new_filename").
		WithExtension("jpg").
		WithPath("new_path")

	mockFuncs := &MockFuncs{}
	mockFuncs.On("FixDuplicates", mock.Anything).Return(file, nil)
	defer helpers.MockFunction(&filerepo.FixDuplicates, mockFuncs.FixDuplicates, func() {})()

	defer helpers.MockFunction(&filerepo.GetFileByID, mockGetFileByID, func() {})()

	expected := *file
	expected.Filename = fvp.Filename
	expected.Extension = fvp.Extension
	expected.Path = fvp.Path

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, expected.GetValidFieldsData())
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	f, err := repo.UpdateFileEntry(context.Background(), file, fvp, owner)
	assert.NoError(t, err)
	assert.Equal(t, expected.Filename, f.Filename)
	assert.Equal(t, expected.Extension, f.Extension)
	assert.Equal(t, expected.Path, f.Path)
}

func TestUpdateFileEntry_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	defer helpers.MockFunction(&filerepo.GetFileByID, mockGetFileByID, func() {})()

	file := makeUserFile()
	owner := makeOwner()
	fvp := fschema.NewUserFileBuilder().
		WithFilename("new_filename").
		WithExtension("jpg").
		WithPath("new_path")

	mockFuncs := &MockFuncs{}
	mockFuncs.On("FixDuplicates", mock.Anything).Return(file, nil)
	release := helpers.MockFunction(&filerepo.FixDuplicates, mockFuncs.FixDuplicates, func() {})
	defer release()

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything).Return(fmt.Errorf("update error"))
	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	_, err := repo.UpdateFileEntry(context.Background(), file, fvp, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
}
