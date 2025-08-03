package globals

import (
	"os"
	"testing"

	"homeserver/test"

	"github.com/stretchr/testify/assert"
)

func TestGlobals(t *testing.T) {
	mockPages := []string{"/login", "/search/", "/"}
	cleanupPages := test.MockGlobalVar(&Pages, mockPages, func() {})
	defer cleanupPages()

	mockProtectedPages := []string{"/settings", "/upload"}
	cleanupProtectedPages := test.MockGlobalVar(&ProtectedPages, mockProtectedPages, func() {})
	defer cleanupProtectedPages()

	mockAvailableFileTypes := []string{"jpg", "png", "txt"}
	cleanupAvailableFileTypes := test.MockGlobalVar(&AvailableFileTypes, mockAvailableFileTypes, func() {})
	defer cleanupAvailableFileTypes()

	t.Run("IsPage returns true for valid pages", func(t *testing.T) {
		assert.True(t, IsPage("/login"))
		assert.True(t, IsPage("/search/"))
		assert.True(t, IsPage("/"))
		assert.False(t, IsPage("/invalid"))
	})

	t.Run("IsProtectedPath returns true for protected pages", func(t *testing.T) {
		assert.True(t, IsProtectedPath("/settings"))
		assert.True(t, IsProtectedPath("/upload"))
		assert.False(t, IsProtectedPath("/login"))
		assert.False(t, IsProtectedPath("/invalid"))
	})

	t.Run("IsFileTypeAvailable returns true for available file types", func(t *testing.T) {
		assert.True(t, IsFileTypeAvailable("jpg"))
		assert.True(t, IsFileTypeAvailable("PNG"))
		assert.False(t, IsFileTypeAvailable("exe"))
	})

	t.Run("BuildFilePath constructs correct file paths", func(t *testing.T) {
		os.Setenv(EnvStoragePath, "/storage/")

		assert.Equal(t, "/storage/relative/filename.txt", BuildFilePath("relative", "filename", "txt"))
		assert.Equal(t, "/storage/filename.txt", BuildFilePath("", "filename", "txt"))

		assert.Panics(t, func() {
			BuildFilePath("/relative", "filename", "txt")
		}, "Relative path must not start or end with a slash")
		assert.Panics(t, func() {
			BuildFilePath("relative/", "filename", "txt")
		}, "Relative path must not start or end with a slash")

		os.Setenv(EnvStoragePath, "/invalid")
		assert.Panics(t, func() {
			BuildFilePath("relative", "filename", "txt")
		}, "Storage path must end with a slash")
	})
}
