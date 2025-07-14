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
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE (username = $1 OR email = $1) AND password = $2 AND email_verified = true`, form.Identificator, form.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (db *Database) CheckToken(ctx context.Context, token string) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE auth_token = $1 AND email_verified = true`, token).Scan(&userID)
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

func (db *Database) QueryUsername(ctx context.Context, username string) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE username = $1`, username).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error querying username: %v", err)
	}
	return userID, nil
}

func (db *Database) QueryEmail(ctx context.Context, email string) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error querying email: %v", err)
	}
	return userID, nil
}

func (db *Database) QueryPassword(ctx context.Context, password string) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE password = $1`, password).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error querying password: %v", err)
	}
	return userID, nil
}

func (db *Database) VerifyToken(ctx context.Context, token string) error {
	var userID int
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE auth_token = $1 AND email_verified = false`, token).Scan(&userID)
	if err != nil {
		return fmt.Errorf("error verifying token: %v", err)
	}
	if userID == 0 {
		return fmt.Errorf("invalid token")
	}

	tag, err := db.pool.Exec(ctx, `UPDATE users SET email_verified = true WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("error updating email verification status: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("affected rows mismatch: expected 1, got %d", tag.RowsAffected())
	}

	// The command tag is returned by Exec and contains information about the executed SQL command,
	// such as the number of rows affected. You might need it if you want to check how many rows were updated,
	// or to verify that the command executed as expected.
	return nil
}
