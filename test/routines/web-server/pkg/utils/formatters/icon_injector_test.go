package formatters_test

import (
	"homeserver/test/helpers"
	"homeserver/web-server/pkg/utils/formatters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIconLoader struct {
	mock.Mock
}

func (m *MockIconLoader) LoadIcon(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func TestIconsInjector(t *testing.T) {
	mockIconLoader := new(MockIconLoader)
	cleanup := helpers.MockFunction(&formatters.LoadIcon, mockIconLoader.LoadIcon, func() {})
	defer cleanup()

	mockIconLoader.On("LoadIcon", "test-icon").Return("<svg>...</svg>")
	mockIconLoader.On("LoadIcon", "another-icon").Panic("Icon not found")

	// Table-driven tests for non-panicking cases
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Injects icons correctly",
			input:    []byte("<icon name=\"test-icon\"></icon> and <icon name=\"test-icon\"></icon>"),
			expected: []byte("<svg>...</svg> and <svg>...</svg>"),
		},
		{
			name:     "Handles empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			name:     "Handles input without icon tags",
			input:    []byte("No icons here."),
			expected: []byte("No icons here."),
		},
		{
			name:     "Handles no name tag but icon name valid",
			input:    []byte("<icon \"test-icon\"></icon>"),
			expected: []byte("<svg>...</svg>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatters.IconsInjector(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Tests for panicking cases
	t.Run("Handles malformed icon tags", func(t *testing.T) {
		input := []byte("<icon name=\"test-icon\">")

		assert.PanicsWithValue(t, formatters.ErrorMalformedTag, func() {
			formatters.IconsInjector(input)
		})
	})

	t.Run("Handles missing name attribute", func(t *testing.T) {
		input := []byte("<icon></icon>")

		assert.PanicsWithValue(t, formatters.ErrorMissingNameAttr, func() {
			formatters.IconsInjector(input)
		})
	})

	t.Run("Handles non-existent icon", func(t *testing.T) {
		input := []byte("<icon name=\"another-icon\"></icon>")

		assert.PanicsWithValue(t, "Icon not found", func() {
			formatters.IconsInjector(input)
		})
	})

	t.Run("Handles malformed name attribute without quotes", func(t *testing.T) {
		input := []byte("<icon name=></icon>")

		assert.PanicsWithValue(t, formatters.ErrorMissingNameAttr, func() {
			formatters.IconsInjector(input)
		})
	})

	t.Run("Handles malformed name attribute with 1 quote", func(t *testing.T) {
		input := []byte("<icon name=\"another-icon></icon>")

		assert.PanicsWithValue(t, formatters.ErrorMalformedTag, func() {
			formatters.IconsInjector(input)
		})
	})

	mockIconLoader.AssertExpectations(t)
}
