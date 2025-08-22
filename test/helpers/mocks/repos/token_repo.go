package repomocks

import (
	"context"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/stretchr/testify/mock"
)

type MockTokenRepo struct {
	mock.Mock
}

func NewMockTokenRepo() *MockTokenRepo {
	return &MockTokenRepo{}
}

func (m *MockTokenRepo) QueryToken(ctx context.Context, token string) (*uschema.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*uschema.User), args.Error(1)
}

func (m *MockTokenRepo) CreateToken(ctx context.Context, userID int, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockTokenRepo) DeleteToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepo) CheckToken(ctx context.Context, token string) (int, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockTokenRepo) VerifyToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
