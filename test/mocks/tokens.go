package mocks

import "github.com/stretchr/testify/mock"

type MockTokens struct {
	mock.Mock
}

func (m *MockTokens) NewToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockTokens) ValidateToken(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}
