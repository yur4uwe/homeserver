package user

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas/user_schema"
)

// CREATE

func CreateUserEntry(ctx context.Context, db interfaces.DatabaseInterface, user *user_schema.User) (*user_schema.User, error) {
	var insertedUser user_schema.User

	err := db.QueryRow(
		ctx,
		`INSERT INTO users (username, password, email, auth_token, created_at, email_verified)
         VALUES ($1, $2, $3, $4, NOW(), $5)
         RETURNING id, username, password, email, auth_token, created_at, email_verified`,
		user.Username,
		user.Password,
		user.Email,
		user.AuthToken,
		user.EmailVerified,
	).Scan(
		&insertedUser.ID,
		&insertedUser.Username,
		&insertedUser.Password,
		&insertedUser.Email,
		&insertedUser.AuthToken,
		&insertedUser.CreatedAt,
		&insertedUser.EmailVerified,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user entry: %v", err)
	}

	return &insertedUser, nil
}

// UPDATE

// DELETE
