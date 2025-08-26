package utils

import (
	"homeserver/constants"
	"homeserver/test/helpers"
	"homeserver/web-server/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailableFileTypes(t *testing.T) {
	mockAvailableFileTypes := []string{"jpg", "png", "txt"}
	cleanupAvailableFileTypes := helpers.MockGlobalVar(&constants.AvailableFileTypes, mockAvailableFileTypes, func() {})
	defer cleanupAvailableFileTypes()

	t.Run("IsFileTypeAvailable returns true for available file types", func(t *testing.T) {
		assert.True(t, utils.IsFileTypeAvailable("jpg"))
		assert.True(t, utils.IsFileTypeAvailable("PNG"))
		assert.False(t, utils.IsFileTypeAvailable("exe"))
	})
}
