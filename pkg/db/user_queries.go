package db

import (
	"context"
	"fmt"
)

func (db *Database) CreateUserEntry(ctx context.Context, user *User) error {
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

func (db *Database) ExtractUserEntries(ctx context.Context) ([]*User, error) {
	rows, err := db.pool.Query(
		ctx,
		`SELECT id, username, password, email, auth_token, created_at, email_verified FROM users`,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying user entries: %v", err)
	}
	defer rows.Close()

	var entries []*User
	for rows.Next() {
		entry := &User{}
		if err := rows.Scan(&entry.ID, &entry.Username, &entry.Password, &entry.Email, &entry.AuthToken, &entry.CreatedAt, &entry.EmailVerified); err != nil {
			return nil, fmt.Errorf("error scanning user entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (db *Database) CheckUserFolder(ctx context.Context, user *User) error {
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
		folder := &UserFolders{
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
