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

func (db *Database) CreateTable(ctx context.Context, table *globals.DBTable) error {
	fmt.Printf("Creating table %s...\n", table.Name)
	_, err := db.pool.Exec(ctx, table.SQLCreation)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", table.Name, err)
	}

	fmt.Printf("Table %s created successfully.\n", table.Name)
	return nil
}

func (db *Database) CreateUser(ctx context.Context, username, email, password string) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`, username, email, password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (db *Database) QueryUserOnLogin(ctx context.Context, form globals.LoginForm) (int, bool, error) {
	var userID int
	var email_verified bool
	err := db.pool.QueryRow(ctx, `SELECT id, email_verified FROM users WHERE (username = $1 OR email = $1) AND password = $2`, form.Identificator, form.Password).
		Scan(&userID, &email_verified)
	if err != nil {
		return 0, false, err
	}
	return userID, email_verified, nil
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
