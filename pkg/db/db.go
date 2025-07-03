package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (db *Database) CreateUser(ctx context.Context, username, email, password string) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`, username, email, password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (db *Database) QueryUserOnLogin(ctx context.Context, form globals.LoginForm) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE (username = $1 OR email = $1) AND password = $2`, form.Identificator, form.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (db *Database) CheckToken(ctx context.Context, token string) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE auth_token = $1`, token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error checking token: %v", err)
	}
	return userID, nil
}

func (db *Database) CreateToken(ctx context.Context, userID int, token string) error {
	_, err := db.pool.Exec(ctx, `UPDATE users SET auth_token = $1 WHERE id = $2`, token, userID)
	if err != nil {
		return fmt.Errorf("error creating token: %v", err)
	}
	return nil
}
