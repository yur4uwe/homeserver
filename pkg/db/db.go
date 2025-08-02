package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/user_schema"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseInterface interface {
	QueryToken(ctx context.Context, token string) (*user_schema.User, error)
}

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase() (*Database, error) {
	dbURL := os.Getenv("DB_URL")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	fmt.Println("Connected to the database successfully.")
	return &Database{pool: pool}, nil
}

func (db *Database) Close() {
	db.pool.Close()
	fmt.Println("Database connection closed.")
}

func (db *Database) TableExists(ctx context.Context, tableName string) (bool, error) {
	var exists bool
	err := db.pool.QueryRow(ctx, `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)`, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if table exists: %v", err)
	}
	return exists, nil
}

func (db *Database) CreateTable(ctx context.Context, table *globals.DBTable) error {
	fmt.Printf("Creating table %s...\n", table.Name)
	_, err := db.pool.Exec(ctx, table.SQLCreation)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", table.Name, err)
	}

	fmt.Printf("Table %s created successfully.\n", table.Name)
	return nil
}
