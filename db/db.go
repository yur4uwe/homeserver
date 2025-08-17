package db

import (
	"context"
	"fmt"
	"homeserver/interfaces"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool interfaces.QueryExecutor
}

func NewConnection() (*Database, error) {
	dbURL := os.Getenv("DB_URL")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	fmt.Println("Connected to the database successfully.")
	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	db.Pool.(*pgxpool.Pool).Close()
}
