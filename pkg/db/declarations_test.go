package db

import (
	"context"
	"fmt"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/test"
	"io"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type cleanupFunc func()

var (
	mockUser           *user_schema.User
	mockSecondUser     *user_schema.User
	mockIncompleteUser *user_schema.User
	mockRequester      *user_schema.User
	mockFolder         *folder_schema.UserFolder
	mockFile           *file_schema.UserFile
	connection         *Database
)

func NewTestDatabase() (*Database, error) {
	testDBUrl := os.Getenv("TEST_DB_URL")

	pool, err := pgxpool.New(context.Background(), testDBUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	fmt.Println("Connected to the database successfully.")
	return &Database{pool: pool}, nil
}

func (db *Database) UseTransaction() cleanupFunc {
	if db == nil {
		panic("Database connection is nil")
	}

	tx, err := db.pool.Begin(context.Background())
	if err != nil {
		return nil
	}

	cleanup := test.MockGlobalVar(&db.pool, tx, func() {
		tx.Rollback(context.Background())
	})

	return cleanup
}

func (db *Database) Copy() *Database {
	if db == nil {
		panic("Database connection is nil")
	}

	return &Database{
		pool: db.pool,
	}
}

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

type MockIO struct {
	mock.Mock
}

func (m *MockIO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	args := m.Called(dst, src)
	return args.Get(0).(int64), args.Error(1)
}
