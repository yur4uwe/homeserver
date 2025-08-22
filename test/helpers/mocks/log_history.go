package mocks

import (
	"github.com/stretchr/testify/mock"
)

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
