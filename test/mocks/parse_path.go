package mocks

import "github.com/stretchr/testify/mock"

type ParsePathMock struct {
	mock.Mock
}

func (m *ParsePathMock) ParsePathWithUName(path string, username string) (string, []string, error) {
	args := m.Called(path, username)
	return args.Get(0).(string), args.Get(1).([]string), args.Error(2)
}
