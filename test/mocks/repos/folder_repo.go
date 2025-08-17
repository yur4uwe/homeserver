package repomocks

import (
	"context"

	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/stretchr/testify/mock"
)

type MockFolderRepo struct {
	mock.Mock
}

func NewMockFolderRepo() *MockFolderRepo {
	return &MockFolderRepo{}
}

func (m *MockFolderRepo) GetFolderByID(ctx context.Context, id int) (*folschema.UserFolder, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*folschema.UserFolder), args.Error(1)
}

func (m *MockFolderRepo) GetFolderByPath(ctx context.Context, folder *folschema.UserFolder) (*folschema.UserFolder, error) {
	args := m.Called(ctx, folder)
	return args.Get(0).(*folschema.UserFolder), args.Error(1)
}

func (m *MockFolderRepo) GetFolderByOwner(ctx context.Context, ownerID int) (*folschema.UserFolder, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(*folschema.UserFolder), args.Error(1)
}

func (m *MockFolderRepo) GetChildren(ctx context.Context, parentID int) ([]*folschema.UserFolder, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]*folschema.UserFolder), args.Error(1)
}

func (m *MockFolderRepo) CreateFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	args := m.Called(ctx, folder, owner)
	return args.Get(0).(*folschema.UserFolder), args.Error(1)
}

func (m *MockFolderRepo) UpdateFolder(ctx context.Context, folder *folschema.UserFolder, fieldsToUpdate []string, owner *uschema.User) error {
	args := m.Called(ctx, folder, fieldsToUpdate, owner)
	return args.Error(0)
}

func (m *MockFolderRepo) DeleteFolderByID(ctx context.Context, id int, ownerID int) (int, error) {
	args := m.Called(ctx, id, ownerID)
	return args.Int(0), args.Error(1)
}

func (m *MockFolderRepo) DeleteFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (int, error) {
	args := m.Called(ctx, folder, owner)
	return args.Int(0), args.Error(1)
}
