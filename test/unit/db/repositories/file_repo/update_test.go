package filerepo_test

import (
	"context"
	"errors"
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

func TestUpdateFileEntry_Success(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := makeUserFile()
	owner := makeOwner()
	fvp := map[string]any{
		"filename":  "new_filename",
		"extension": "jpg",
		"path":      "new_path",
	}

	mockFuncs := &MockFuncs{}
	mockFuncs.On("FixDuplicates", mock.Anything).Return(file, nil)
	release := helpers.MockFunction(&filerepo.FixDuplicates, mockFuncs.FixDuplicates, func() {})
	defer release()

	expected := *file
	expected.ID += 1
	expected.Filename = fvp["filename"].(string)
	expected.Extension = fvp["extension"].(string)
	expected.Path = fvp["path"].(string)

	mockRows := &mocks.MockRows{}
	mockRows.On("FieldDescriptions").Return(fieldDescs)
	mockRows.On("Close").Return().Once()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	mockRows.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).([]any)
		mocks.RunScan(dest, filerepo.OrderFileUpdateFields(fschema.FieldsList(), &expected))
	})

	mockPool.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

	err := repo.UpdateFileEntry(context.Background(), file, fvp, owner)
	assert.NoError(t, err)
	assert.NotEqual(t, file.ID, 1)
	assert.Equal(t, expected.Filename, file.Filename)
	assert.Equal(t, expected.Extension, file.Extension)
	assert.Equal(t, expected.Path, file.Path)
}

func TestUpdateFileEntry_PreventDuplicateError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := makeUserFile()
	owner := makeOwner()
	fvp := map[string]any{
		"filename":  "new_filename",
		"extension": "jpg",
		"path":      "new_path",
	}

	// Patch FixDuplicates to return error
	mockFuncs := &MockFuncs{}
	mockFuncs.On("FixDuplicates", mock.Anything).Return((*fschema.UserFile)(nil), errors.New("duplicate error"))
	release := helpers.MockFunction(&filerepo.FixDuplicates, mockFuncs.FixDuplicates, func() {})
	defer release()

	err := repo.UpdateFileEntry(context.Background(), file, fvp, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preventing duplicate file name")
}

func TestUpdateFileEntry_DBError(t *testing.T) {
	mockPool := &mocks.MockQueryExecutor{}
	mockDB := &db.Database{Pool: mockPool}
	repo := filerepo.NewFileRepo(mockDB)

	file := makeUserFile()
	owner := makeOwner()
	fvp := map[string]any{
		"filename":  "new_filename",
		"extension": "jpg",
		"path":      "new_path",
	}

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

	err := repo.UpdateFileEntry(context.Background(), file, fvp, owner)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
}
