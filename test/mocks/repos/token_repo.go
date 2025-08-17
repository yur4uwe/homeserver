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
