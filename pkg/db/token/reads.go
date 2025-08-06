package token

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas/user_schema"
	"os"
)

var (
	Mkdir = os.Mkdir
)

func CheckToken(ctx context.Context, db interfaces.DatabaseInterface, token string) (int, error) {
	var userID int
	err := db.QueryRow(ctx, `SELECT id FROM users WHERE auth_token = $1 AND email_verified = true`, token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error checking token: %v", err)
	}
	return userID, nil
}

func QueryToken(ctx context.Context, db interfaces.DatabaseInterface, token string) (*user_schema.User, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	var user user_schema.User
	err := db.QueryRow(ctx, `
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

func VerifyToken(ctx context.Context, db interfaces.DatabaseInterface, token string) error {
	var userName string
	err := db.QueryRow(ctx, `SELECT username FROM users WHERE auth_token = $1 AND email_verified = false`, token).Scan(&userName)
	if err != nil {
		return fmt.Errorf("error verifying token: %v", err)
	}
	if userName == "" {
		return fmt.Errorf("invalid token")
	}

	fmt.Println("Token verified for user:", userName)

	tag, err := db.Exec(ctx, `UPDATE users SET email_verified = true WHERE username = $1`, userName)
	if err != nil {
		return fmt.Errorf("error updating email verification status: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("affected rows mismatch: expected 1, got %d", tag.RowsAffected())
	}

	storagePath := os.Getenv(globals.EnvStoragePath)

	if err := Mkdir(storagePath+userName, 0755); err != nil {
		return fmt.Errorf("error creating user directory: %v", err)
	}

	return nil
}
