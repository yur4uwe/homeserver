package middleware

import (
	"context"
	"html/template"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockParseFiles struct {
	mock.Mock
}

func (m *MockParseFiles) ParseFiles(filenames ...string) (*template.Template, error) {
	args := m.Called(filenames)
	return args.Get(0).(*template.Template), args.Error(1)
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

func (m *MockDB) QueryToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
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

type mockParseAsString struct {
	mock.Mock
}

func (m *mockParseAsString) ParseTemplateToString(templatePath string, data interface{}) (string, error) {
	args := m.Called(templatePath, data)
	return args.String(0), args.Error(1)
}

type MockErrorResponse struct {
	mock.Mock
}

func (m *MockErrorResponse) ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	args := m.Called(w, message, statusCode)
	if args.Get(0) != nil {
		http.Error(w, message, statusCode)
	}
}
