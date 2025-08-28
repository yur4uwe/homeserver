package db

import (
	"context"
	"fmt"
	"homeserver/db/repositories/helpers"
	"homeserver/interfaces"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database[T any] struct {
	Pool interfaces.QueryExecutor
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

func (db *Database[T]) Close() {
	db.Pool.(*pgxpool.Pool).Close()
}

func (db *Database[T]) QueryOne(ctx context.Context, query string, args ...any) (T, error) {
	var zero T
	rows, err := db.QueryMany(ctx, query, args...)
	if err != nil {
		return zero, err
	}
	if len(rows) != 1 {
		return zero, fmt.Errorf("unexpected error: %d rows returned", len(rows))
	}
	return rows[0], nil
}

func (db *Database[T]) QueryMany(ctx context.Context, query string, args ...any) ([]T, error) {
	dest := []T{}

	rows, err := db.Pool.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	if err = helpers.ScanRows(rows, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}
