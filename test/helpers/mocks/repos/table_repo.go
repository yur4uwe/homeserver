package repomocks

import (
	"context"
	"homeserver/constants/dbconsts"

	"github.com/stretchr/testify/mock"
)

type MockTableRepo struct {
	mock.Mock
}

func NewMockTableRepo() *MockTableRepo {
	return &MockTableRepo{}
}

func (m *MockTableRepo) DropTable(ctx context.Context, table_name string) error {
	args := m.Called(ctx, table_name)
	return args.Error(0)
}

func (m *MockTableRepo) TableExists(ctx context.Context, table_name string) (bool, error) {
	args := m.Called(ctx, table_name)
	return args.Bool(0), args.Error(1)
}

func (m *MockTableRepo) CreateTable(ctx context.Context, table *dbconsts.DBTable) error {
	args := m.Called(ctx, table)
	return args.Error(0)
}
