package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockIO struct {
	mock.Mock
}

func (m *MockIO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	args := m.Called(dst, src)
	return args.Get(0).(int64), args.Error(1)
}
