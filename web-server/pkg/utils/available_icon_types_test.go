package utils

import (
	"homeserver/constants"
	"homeserver/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailableFileTypes(t *testing.T) {
	mockAvailableFileTypes := []string{"jpg", "png", "txt"}
	cleanupAvailableFileTypes := test.MockGlobalVar(&constants.AvailableFileTypes, mockAvailableFileTypes, func() {})
	defer cleanupAvailableFileTypes()

	t.Run("IsFileTypeAvailable returns true for available file types", func(t *testing.T) {
		assert.True(t, IsFileTypeAvailable("jpg"))
		assert.True(t, IsFileTypeAvailable("PNG"))
		assert.False(t, IsFileTypeAvailable("exe"))
	})
}
