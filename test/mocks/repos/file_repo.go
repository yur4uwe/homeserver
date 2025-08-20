package repomocks

import (
	"context"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/stretchr/testify/mock"
)

type MockFileRepo struct {
	mock.Mock
}

func NewMockFileRepo() *MockFileRepo {
	return &MockFileRepo{}
}

func (m *MockFileRepo) DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (int, error) {
	args := m.Called(ctx, file, owner)
	return args.Int(0), args.Error(1)
}

func (m *MockFileRepo) UpdateFileEntry(ctx context.Context, file *fschema.UserFile, fvp map[string]any, owner *uschema.User) error {
	args := m.Called(ctx, file, fvp, owner)
	return args.Error(0)
}

func (m *MockFileRepo) CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	args := m.Called(ctx, file, owner)
	return args.Get(0).(*fschema.UserFile), args.Error(1)
}

func (m *MockFileRepo) GetFileByID(ctx context.Context, id string) (*fschema.UserFile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*fschema.UserFile), args.Error(1)
}

func (m *MockFileRepo) GetFullFileInfo(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	args := m.Called(ctx, file, owner)
	return args.Get(0).(*fschema.UserFile), args.Error(1)
}

func (m *MockFileRepo) GetFilesByUser(ctx context.Context, owner *uschema.User) ([]fschema.UserFile, error) {
	args := m.Called(ctx, owner)
	return args.Get(0).([]fschema.UserFile), args.Error(1)
}

func (m *MockFileRepo) FuzzySearch(ctx context.Context, search_query string, directory string, owner *uschema.User) ([]fschema.UserFile, error) {
	args := m.Called(ctx, search_query, directory, owner)
	return args.Get(0).([]fschema.UserFile), args.Error(1)
}

func (m *MockFileRepo) GetFilesByDirectory(ctx context.Context, directory string, owner *uschema.User) ([]fschema.UserFile, error) {
	args := m.Called(ctx, directory, owner)
	return args.Get(0).([]fschema.UserFile), args.Error(1)
}

func (m *MockFileRepo) DeleteFileByID(ctx context.Context, fileID string, owner *uschema.User) error {
	args := m.Called(ctx, fileID, owner)
	return args.Error(0)
}
