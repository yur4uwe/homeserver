package user

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
)

func ExtractUserEntries(ctx context.Context, db interfaces.DatabaseInterface) ([]*user_schema.User, error) {
	rows, err := db.Query(
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

func CheckUserFolder(ctx context.Context, db interfaces.DatabaseInterface, user *user_schema.User) error {
	var exists bool
	err := db.QueryRow(
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

func QueryUserOnLogin(ctx context.Context, db interfaces.DatabaseInterface, form globals.LoginForm) (int, bool, error) {
	var userID int
	var email_verified bool
	err := db.QueryRow(ctx, `SELECT id, email_verified FROM users WHERE (username = $1 OR email = $1) AND password = $2`, form.Identificator, form.Password).
		Scan(&userID, &email_verified)
	if err != nil {
		return 0, false, err
	}
	return userID, email_verified, nil
}

func QueryUsername(ctx context.Context, db interfaces.DatabaseInterface, username string) (int, error) {
	var userID int
	err := db.QueryRow(ctx, `SELECT id FROM users WHERE username = $1`, username).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error querying username: %v", err)
	}
	return userID, nil
}

func QueryEmail(ctx context.Context, db interfaces.DatabaseInterface, email string) (int, error) {
	var userID int
	err := db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error querying email: %v", err)
	}
	return userID, nil
}

func QueryPassword(ctx context.Context, db interfaces.DatabaseInterface, password string) (int, error) {
	var userID int
	err := db.QueryRow(ctx, `SELECT id FROM users WHERE password = $1`, password).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error querying password: %v", err)
	}
	return userID, nil
}
