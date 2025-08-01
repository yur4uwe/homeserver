package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
)

func (db *Database) CreateUserEntry(ctx context.Context, user *user_schema.User) error {
	_, err := db.pool.Exec(
		ctx,
		`INSERT INTO users (username, password, email, auth_token, created_at, email_verified) VALUES ($1, $2, $3, $4, NOW(), $5)`,
		user.Username,
		user.Password,
		user.Email,
		user.AuthToken,
		user.EmailVerified,
	)
	if err != nil {
		return fmt.Errorf("error creating user entry: %v", err)
	}
	return nil
}

func (db *Database) ExtractUserEntries(ctx context.Context) ([]*user_schema.User, error) {
	rows, err := db.pool.Query(
		ctx,
		`SELECT id, username, password, email, auth_token, created_at, email_verified FROM users`,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying user entries: %v", err)
	}
	defer rows.Close()

	var entries []*user_schema.User
	for rows.Next() {
		entry := &user_schema.User{}
		if err := rows.Scan(&entry.ID, &entry.Username, &entry.Password, &entry.Email, &entry.AuthToken, &entry.CreatedAt, &entry.EmailVerified); err != nil {
			return nil, fmt.Errorf("error scanning user entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (db *Database) CheckUserFolder(ctx context.Context, user *user_schema.User) error {
	var exists bool
	err := db.pool.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT id FROM user_folders WHERE owner_id = $1)`,
		user.ID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user folder: %v", err)
	}

	if !exists {
		folder := &folder_schema.UserFolder{
			Name:    user.Username,
			Parent:  nil,
			OwnerID: user.ID,
		}

		if err := db.CreateFolderEntry(ctx, folder); err != nil {
			return fmt.Errorf("error creating folder for user %s: %v", user.Username, err)
		}
	}

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
