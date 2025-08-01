package folder_schema

import (
	"homeserver/pkg/globals"
	"homeserver/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGetEnv struct {
	mock.Mock
}

func (m *MockGetEnv) GetEnv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestFolderBuilderSuite(t *testing.T) {
	mockGetEnv := new(MockGetEnv)
	mockGetEnv.On("GetEnv", globals.EnvStoragePath).Return("storage/")

	cleanup := test_utils.MockFunction(&getStoragePath, mockGetEnv.GetEnv, func() {})
	defer cleanup()

	t.Run("TestNewUserFolderBuilder", func(t *testing.T) {
		newFolder := NewUserFolderBuilder()

		assert.NotNil(t, newFolder, "New folder should not be nil")
		assert.Equal(t, defaultFolder, *newFolder, "New folder should match the default folder")
	})

	t.Run("TestWithID", func(t *testing.T) {
		newFolder := NewUserFolderBuilder().WithID(1)
		assert.Equal(t, 1, newFolder.ID, "Folder ID should be set to 1")
	})

	t.Run("TestWithIDPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFolderBuilder().WithID(0)
		}, "Setting folder ID to 0 should panic")
	})

	t.Run("TestWithName", func(t *testing.T) {
		newFolder := NewUserFolderBuilder().WithName("documents")
		assert.Equal(t, "documents", newFolder.Name, "Folder name should be set to 'documents'")
	})

	t.Run("TestWithNamePanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFolderBuilder().WithName("")
		}, "Setting an empty folder name should panic")
	})

	t.Run("TestWithPath", func(t *testing.T) {
		newFolder := NewUserFolderBuilder().WithPath("john_doe")
		assert.Equal(t, "john_doe", newFolder.Path, "Folder path should be set to 'john_doe'")
	})

	t.Run("TestWithPathPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFolderBuilder().WithPath("")
		}, "Setting an empty folder path should panic")

		assert.Panics(t, func() {
			NewUserFolderBuilder().WithPath("/invalid_path/")
		}, "Setting a folder path starting or ending with a slash should panic")
	})

	t.Run("TestWithParent", func(t *testing.T) {
		parentID := 42
		newFolder := NewUserFolderBuilder().WithParent(&parentID)
		assert.Equal(t, &parentID, newFolder.Parent, "Parent folder ID should be set to 42")
	})

	t.Run("TestWithParentPanic", func(t *testing.T) {
		parentID := 0
		assert.Panics(t, func() {
			NewUserFolderBuilder().WithParent(&parentID)
		}, "Setting parent folder ID to 0 should panic")
	})

	t.Run("TestWithOwnerID", func(t *testing.T) {
		newFolder := NewUserFolderBuilder().WithOwnerID(100)
		assert.Equal(t, 100, newFolder.OwnerID, "Owner ID should be set to 100")
	})

	t.Run("TestWithOwnerIDPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserFolderBuilder().WithOwnerID(0)
		}, "Setting owner ID to 0 should panic")
	})

	t.Run("TestAbsolutePath", func(t *testing.T) {
		newFolder := NewUserFolderBuilder().
			WithName("documents").
			WithPath("john_doe")

		expectedPath := "storage/john_doe/documents"
		assert.Equal(t, expectedPath, newFolder.AbsolutePath(), "Absolute path should match the expected path")
	})
}
