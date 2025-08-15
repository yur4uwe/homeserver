package mocks

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) Close() {
	m.Called()
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	rows_affected := args.Get(0).(int)
	return pgconn.NewCommandTag(fmt.Sprintf("SELECT %d", rows_affected))
}

// --- Add these methods to complete the interface ---

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRows) Values() ([]any, error) {
	args := m.Called()
	return args.Get(0).([]any), args.Error(1)
}

func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([][]byte)
}

func (m *MockRows) Conn() *pgx.Conn {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgx.Conn)
}
