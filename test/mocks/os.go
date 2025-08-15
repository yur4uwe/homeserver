package mocks

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockOS struct {
	mock.Mock
}

func (m *MockOS) Mkdir(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockOS) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockOS) Rename(oldPath, newPath string) error {
	args := m.Called(oldPath, newPath)
	return args.Error(0)
}

func (m *MockOS) RemoveAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockOS) Create(name string) (*os.File, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockOS) Remove(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockOS) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockOS) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockOS) ReadDir(name string) ([]os.DirEntry, error) {
	args := m.Called(name)
	return args.Get(0).([]os.DirEntry), args.Error(1)
}

func (m *MockOS) WriteFile(name string, data []byte, perm os.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

func (m *MockOS) Getenv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

type MockDirEntry struct {
	mock.Mock
}

func (m *MockDirEntry) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockDirEntry) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockDirEntry) Type() os.FileMode {
	args := m.Called()
	return args.Get(0).(os.FileMode)
}

func (m *MockDirEntry) Info() (os.FileInfo, error) {
	args := m.Called()
	return args.Get(0).(os.FileInfo), args.Error(1)
}
