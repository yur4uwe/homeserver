package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockQueryExecutor struct {
	mock.Mock
}

// Basic Query Methods
func (m *MockQueryExecutor) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Rows), argsMock.Error(1)
}

func (m *MockQueryExecutor) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Row)
}

func (m *MockQueryExecutor) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgconn.CommandTag), argsMock.Error(1)
}
