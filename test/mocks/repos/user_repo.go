package repomocks

import "github.com/stretchr/testify/mock"

type MockUserRepo struct {
	mock.Mock
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{}
}
