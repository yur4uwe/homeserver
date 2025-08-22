package utils

import (
	"homeserver/constants"
	"homeserver/test/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageUtils(t *testing.T) {
	mockPages := []string{"/login", "/search/", "/"}
	cleanupPages := helpers.MockGlobalVar(&constants.Pages, mockPages, func() {})
	defer cleanupPages()

	mockProtectedPages := []string{"/settings", "/upload"}
	cleanupProtectedPages := helpers.MockGlobalVar(&constants.ProtectedPages, mockProtectedPages, func() {})
	defer cleanupProtectedPages()

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
}
