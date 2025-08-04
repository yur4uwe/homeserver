package file

import (
	"homeserver/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOS struct {
	mock.Mock
}

func (m *MockOS) Stat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func TestPreventDuplicateFileName(t *testing.T) {
	os.Setenv("STORAGE_PATH", "storage/") // Set environment variable for testing

	mockOS := new(MockOS)

	cleanupStat := test.MockFunction(&Stat, mockOS.Stat, func() {})
	defer cleanupStat()

	// Test cases
	t.Run("Returns original path if file does not exist", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/file.txt").Return(nil, os.ErrNotExist).Once() // Allow two calls

		path := "path/to/file.txt"
		newPath, err := PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.Equal(t, "storage/"+path, newPath, "Path should match the original path")
	})

	t.Run("Generates new path if file exists", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/file.txt").Return(nil, nil).Once()              // Allow two calls
		mockOS.On("Stat", "storage/path/to/file_1.txt").Return(nil, os.ErrNotExist).Once() // File exists

		path := "path/to/file.txt"
		newPath, err := PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.NotEqual(t, path, newPath, "Path should be updated to avoid duplication")
		assert.Equal(t, "storage/path/to/file_1.txt", newPath, "Generated path should match expected value")
	})

	t.Run("Handles multiple duplicates correctly", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/file.txt").Return(nil, nil).Once()              // File exists
		mockOS.On("Stat", "storage/path/to/file_1.txt").Return(nil, nil).Once()            // File exists
		mockOS.On("Stat", "storage/path/to/file_2.txt").Return(nil, os.ErrNotExist).Once() // File does not exist

		path := "path/to/file.txt"
		newPath, err := PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.Equal(t, "storage/path/to/file_2.txt", newPath, "Generated path should match expected value after handling multiple duplicates")
	})

	t.Run("Panics for invalid file path", func(t *testing.T) {
		assert.PanicsWithValue(t, "Invalid file path: ", func() {
			PreventDuplicateFileName("")
		}, "PreventDuplicateFileName should panic for invalid file path")
	})

	t.Run("Handles paths with special characters", func(t *testing.T) {
		mockOS.On("Stat", "storage/path/to/special@file.txt").Return(nil, os.ErrNotExist).Once() // File does not exist

		path := "path/to/special@file.txt"
		newPath, err := PreventDuplicateFileName(path)
		assert.NoError(t, err, "PreventDuplicateFileName should not return an error")
		assert.Equal(t, "storage/"+path, newPath, "Path should match the original path")
	})

	t.Run("Handles paths with no extension", func(t *testing.T) {
		path := "path/to/file"
		assert.Panics(t, func() {
			PreventDuplicateFileName(path)
		}, "PreventDuplicateFileName should panic for invalid file path")
	})
}
