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

func (m *MockFileRepo) DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) error {
	args := m.Called(ctx, file, owner)
	return args.Error(0)
}

func (m *MockFileRepo) UpdateFileEntry(ctx context.Context, modified_file *fschema.UserFile, fields_to_update []string, owner *uschema.User) error {
	args := m.Called(ctx, modified_file, fields_to_update, owner)
	return args.Error(0)
}

func (m *MockFileRepo) CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	args := m.Called(ctx, file, owner)
	return args.Get(0).(*fschema.UserFile), args.Error(1)
}
