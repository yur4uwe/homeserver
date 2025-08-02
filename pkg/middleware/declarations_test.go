package middleware

import (
	"context"
	"homeserver/pkg/schemas/user_schema"
	"html/template"
	"net/http"

	"github.com/stretchr/testify/mock"
)

var (
	mockTemplate      = template.Must(template.New("").Parse("<html><body>{{.Content}}</body></html>"))
	mockIndexTemplate = template.Must(template.New("index.html").Parse("<html><body>{{.Content}}</body></html>"))
)

type MockParseFiles struct {
	mock.Mock
}

func (m *MockParseFiles) ParseFiles(filenames ...string) *template.Template {
	args := m.Called(filenames)
	return args.Get(0).(*template.Template)
}

type MockLogHistory struct {
	mock.Mock
}

func (m *MockLogHistory) Log(entry string) {
	m.Called(entry)
}

func (m *MockLogHistory) ErrLog(entry, at string) {
	m.Called(entry, at)
}

func (m *MockLogHistory) UnloadLogs(header string) {
	m.Called(header)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryToken(ctx context.Context, token string) (*user_schema.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_schema.User), args.Error(1)
}

type MockMust struct {
	mock.Mock
}

func (m *MockMust) Must(value interface{}, err error) interface{} {
	args := m.Called(value, err)
	if args.Get(1) != nil {
		panic(args.Get(1).(error))
	}
	return args.Get(0)
}

type MockParseAsString struct {
	mock.Mock
}

func (m *MockParseAsString) ParseTemplateToString(templatePath string, data interface{}) (string, error) {
	args := m.Called(templatePath, data)
	return args.String(0), args.Error(1)
}

type MockErrorResponse struct {
	mock.Mock
}

func (m *MockErrorResponse) ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	m.Called(w, message, statusCode)

	http.Error(w, message, statusCode)
}

type MockIsProtectedPath struct {
	mock.Mock
}

func (m *MockIsProtectedPath) IsProtectedPath(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}
