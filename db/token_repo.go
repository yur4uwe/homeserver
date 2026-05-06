package db

import (
	"context"
	"fmt"
	"homeserver/constants"
	"os"
)

var (
	Mkdir  = os.MkdirAll
	Getenv = os.Getenv
)

type TokenRepo interface {
	QueryToken(ctx context.Context, token string, shouldBeVerified bool) (*User, error)
	CreateToken(ctx context.Context, userID int, token string) error
	DeleteToken(ctx context.Context, token string) error
	VerifyToken(ctx context.Context, token string) error
}

type pgTokenRepo struct {
	db QueryExecutor
}

func NewTokenRepo(db QueryExecutor) TokenRepo {
	return &pgTokenRepo{db: db}
}

func (r *pgTokenRepo) CreateToken(ctx context.Context, userID int, token string) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if token == "" {
		return fmt.Errorf("empty token")
	}

	tag, err := r.db.Exec(ctx, `UPDATE users SET auth_token = $1 WHERE id = $2`, token, userID)
	if err != nil {
		return fmt.Errorf("error creating token: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}

func (r *pgTokenRepo) QueryToken(ctx context.Context, token string, shouldBeVerified bool) (*User, error) {
	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	query := `
		SELECT id, username, password, email, auth_token, created_at, email_verified 
		FROM users WHERE auth_token = $1 AND email_verified = $2`

	var u User
	err := r.db.QueryRow(ctx, query, token, shouldBeVerified).
		Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.AuthToken, &u.CreatedAt, &u.EmailVerified)

	if err != nil {
		return nil, fmt.Errorf("error querying token: %w", err)
	}

	return &u, nil
}

func (r *pgTokenRepo) VerifyToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("empty token")
	}

	// 1. Get user details
	var u User
	err := r.db.QueryRow(ctx, `SELECT id, username FROM users WHERE auth_token = $1 AND email_verified = false`, token).
		Scan(&u.ID, &u.Username)
	if err != nil {
		return fmt.Errorf("error finding user for token verification: %w", err)
	}

	// 2. Update verification status
	tag, err := r.db.Exec(ctx, `UPDATE users SET email_verified = true WHERE id = $1`, u.ID)
	if err != nil {
		return fmt.Errorf("error updating email verification status: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("affected rows mismatch: expected 1, got %d", tag.RowsAffected())
	}

	// 3. Create storage directory (Legacy logic, should be moved later)
	storagePath := Getenv(constants.StoragePath)
	if err := Mkdir(storagePath+u.Username, 0755); err != nil {
		return fmt.Errorf("error creating user directory: %w", err)
	}

	return nil
}

func (r *pgTokenRepo) DeleteToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("empty token")
	}

	tag, err := r.db.Exec(ctx, `UPDATE users SET auth_token = NULL WHERE auth_token = $1`, token)
	if err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no token found with value %s", token)
	}

	return nil
}
