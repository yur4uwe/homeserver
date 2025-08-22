package mocks

import "github.com/stretchr/testify/mock"

type ParsePathMock struct {
	mock.Mock
}

func (m *ParsePathMock) ParsePathWithUName(path string, username string) ([]string, error) {
	args := m.Called(path, username)
	return args.Get(0).([]string), args.Error(1)
}
