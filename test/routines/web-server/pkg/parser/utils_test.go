package parser_test

import (
	"errors"
	"homeserver/constants"
	"homeserver/test/helpers"
	"homeserver/test/helpers/mocks/mocksos"
	"homeserver/web-server/pkg/parser"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTemplateParseFiles struct {
	mock.Mock
}

func (m *MockTemplateParseFiles) ParseFiles(filenames ...string) (*template.Template, error) {
	args := m.Called(filenames)
	return args.Get(0).(*template.Template), args.Error(1)
}

func TestLoadIcon(t *testing.T) {
	mockOS := new(mocksos.MockOS)
	mockOS.On("ReadFile", constants.Rootd+"static/assets/test-icon.svg").Return([]byte("<svg>Mock Icon</svg>"), nil)

	// Replace osReadFile with the mock
	cleanup := helpers.MockFunction(&parser.OsReadFile, mockOS.ReadFile, func() {})
	defer cleanup()

	t.Run("TestLoadIconSuccess", func(t *testing.T) {
		icon := parser.LoadIcon("test-icon")
		assert.Equal(t, "<svg>Mock Icon</svg>", icon, "Icon content should match the mock data")
	})

	t.Run("TestLoadIconFailure", func(t *testing.T) {
		mockOS.On("ReadFile", constants.Rootd+"static/assets/missing-icon.svg").Return(nil, errors.New("file not found"))

		assert.Panics(t, func() {
			parser.LoadIcon("missing-icon")
		}, "LoadIcon should panic when the file is not found")
	})
}

func TestNewAlert(t *testing.T) {
	mockTemplate := new(MockTemplateParseFiles)
	mockTemplate.On("ParseFiles", []string{constants.Rootd + "static/templates/modal.html"}).Return(template.Must(template.New("mock").Parse("Mock Template: {{.Title}} - {{.Message}}")), nil)

	// Replace templateParseFiles with the mock
	cleanup := helpers.MockFunction(&parser.ParseFiles, mockTemplate.ParseFiles, func() {})
	defer cleanup()

	t.Run("TestNewAlertSuccess", func(t *testing.T) {
		alert := parser.NewAlert("Test Title", "test-icon", "Test Message", "success")
		assert.Equal(t, "Mock Template: Test Title - Test Message", alert, "Alert content should match the mock template")
	})

	t.Run("TestNewAlertTemplateFailure", func(t *testing.T) {
		mockTemplate.On("ParseFiles", []string{constants.Rootd + "static/templates/missing-template.html"}).Return(nil, errors.New("template not found"))

		assert.Panics(t, func() {
			parser.ParseFiles(constants.Rootd + "static/templates/missing-template.html")
		}, "NewAlert should panic when the template file is not found")
	})
}
