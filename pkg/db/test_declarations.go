package db

import (
	"context"
	"fmt"
	"homeserver/test"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type cleanupFunc func()

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

	cleanup := test.MockGlobalVar(&db.pool, tx, func() {})

	return cleanup
}
