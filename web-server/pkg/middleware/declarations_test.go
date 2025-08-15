package middleware

import (
	"html/template"

	"github.com/stretchr/testify/mock"
)

var (
	mockedTemplate    = template.Must(template.New("").Parse("<html><body>{{.Content}}</body></html>"))
	mockIndexTemplate = template.Must(template.New("index.html").Parse("<html><body>{{.Content}}</body></html>"))
)

type MockParseAsString struct {
	mock.Mock
}

func (m *MockParseAsString) ParseTemplateToString(templatePath string, data interface{}) (string, error) {
	args := m.Called(templatePath, data)
	return args.String(0), args.Error(1)
}

type MockIsProtectedPath struct {
	mock.Mock
}

func (m *MockIsProtectedPath) IsProtectedPath(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}
