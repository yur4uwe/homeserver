package generics_test

import (
	"homeserver/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Run("returns value on success", func(t *testing.T) {
		value := generics.Must("test", nil)

		assert.Equal(t, "test", value.(string))
	})

	t.Run("panics on error", func(t *testing.T) {
		assert.PanicsWithValue(t, assert.AnError, func() {
			generics.Must(nil, assert.AnError)
		}, "Must should panic with the error message")
	})
}
