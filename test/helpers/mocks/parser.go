package mocks

import "github.com/stretchr/testify/mock"

type MockParser struct {
	mock.Mock
}

func (m *MockParser) NewAlert(title string, icon string, message string, alertType string) string {
	args := m.Called(title, icon, message, alertType)
	return args.String(0)
}

func (m *MockParser) LoadIcon(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *MockParser) ParseTemplateToString(templateName string, data interface{}) (string, error) {
	args := m.Called(templateName, data)
	return args.String(0), args.Error(1)
}
