package utils

import (
	"homeserver/constants"
	"homeserver/test/helpers"
	"homeserver/web-server/pkg/utils"
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
		assert.True(t, utils.IsPage("/login"))
		assert.True(t, utils.IsPage("/search/"))
		assert.True(t, utils.IsPage("/"))
		assert.False(t, utils.IsPage("/invalid"))
	})

	t.Run("IsProtectedPath returns true for protected pages", func(t *testing.T) {
		assert.True(t, utils.IsProtectedPath("/settings"))
		assert.True(t, utils.IsProtectedPath("/upload"))
		assert.False(t, utils.IsProtectedPath("/login"))
		assert.False(t, utils.IsProtectedPath("/invalid"))
	})
}
