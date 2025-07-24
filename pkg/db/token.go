package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"os"
)

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

func (db *Database) QueryToken(ctx context.Context, token string) (*User, error) {
	var user User
	err := db.pool.QueryRow(ctx, `
        SELECT id, username, password, email, auth_token, created_at, email_verified
        FROM users
        WHERE auth_token = $1
    `, token).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.AuthToken,
		&user.CreatedAt,
		&user.EmailVerified,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying token: %v", err)
	}
	return &user, nil
}

func (db *Database) VerifyToken(ctx context.Context, token string) error {
	var userName string
	err := db.pool.QueryRow(ctx, `SELECT id FROM users WHERE auth_token = $1 AND email_verified = false`, token).Scan(&userName)
	if err != nil {
		return fmt.Errorf("error verifying token: %v", err)
	}
	if userName == "" {
		return fmt.Errorf("invalid token")
	}

	tag, err := db.pool.Exec(ctx, `UPDATE users SET email_verified = true WHERE username = $1`, userName)
	if err != nil {
		return fmt.Errorf("error updating email verification status: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("affected rows mismatch: expected 1, got %d", tag.RowsAffected())
	}

	storagePath := os.Getenv(globals.EnvStoragePath)

	if err := os.Mkdir(storagePath+userName, 0755); err != nil {
		return fmt.Errorf("error creating user directory: %v", err)
	}

	// The command tag is returned by Exec and contains information about the executed SQL command,
	// such as the number of rows affected. You might need it if you want to check how many rows were updated,
	// or to verify that the command executed as expected.
	return nil
}

func (db *Database) DeleteToken(ctx context.Context, token string) error {
	_, err := db.pool.Exec(ctx, `UPDATE users SET auth_token = NULL WHERE auth_token = $1`, token)
	if err != nil {
		return fmt.Errorf("error deleting token: %v", err)
	}
	return nil
}
