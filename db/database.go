package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryExecutor interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

// Database is the "God object" for data access
type Database struct {
	Users   UserRepo
	Tokens  TokenRepo
	Files   FileRepo
	Folders FolderRepo

	Raw QueryExecutor
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		Users:   NewUserRepo(pool),
		Tokens:  NewTokenRepo(pool),
		Files:   NewFileRepo(pool),
		Folders: NewFolderRepo(pool),
		Raw:     pool,
	}
}

func NewPool() *pgxpool.Pool {
	dbURL := os.Getenv("DB_URL")

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	fmt.Println("Connected to the database successfully.")
	return pool
}
