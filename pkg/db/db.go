package db

import (
	"context"
	"fmt"
	"homeserver/pkg/db/tables"
	"homeserver/pkg/globals"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Database struct {
	pool QueryExecutor
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
	db.pool.(*pgxpool.Pool).Close()
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
	_, err := db.pool.Exec(ctx, table.SQLCreation)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", table.Name, err)
	}

	return nil
}

func (db *Database) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, sql, args...)
}

func (db *Database) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

func (db *Database) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

func (db *Database) ClearEntries(ctx context.Context) (string, error) {
	return tables.ClearEntries(ctx, db)
}

func (db *Database) DropTable(ctx context.Context, tableName string) error {
	return tables.DropTable(ctx, db, tableName)
}
