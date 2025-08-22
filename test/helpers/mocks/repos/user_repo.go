package repomocks

import (
	"context"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{}
}

func (m *MockUserRepo) Create(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*uschema.User), args.Error(1)
}

func (m *MockUserRepo) Query(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*uschema.User), args.Error(1)
}

func (m *MockUserRepo) QueryAll(ctx context.Context) ([]uschema.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]uschema.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *uschema.User, field_vals map[string]any) error {
	args := m.Called(ctx, user, field_vals)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, user *uschema.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
