package db

import (
	"context"
	"fmt"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/test"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type cleanupFunc func()

var (
	mockUser           *user_schema.User
	mockSecondUser     *user_schema.User
	mockIncompleteUser *user_schema.User
	mockRequester      *user_schema.User
	mockFolder         *folder_schema.UserFolder
	mockRootFolder     *folder_schema.UserFolder
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
