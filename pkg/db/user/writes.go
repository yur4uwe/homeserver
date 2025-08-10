package user

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas"
	"homeserver/pkg/schemas/user_schema"
)

// CREATE

func CreateUserEntry(ctx context.Context, db interfaces.DatabaseInterface, user *user_schema.User) (*user_schema.User, error) {
	if err := user.ShouldHaveFields(schemas.Username, schemas.Password, schemas.Email); err != nil {
		return nil, fmt.Errorf("invalid user data: %v", err)
	}

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

func ModifyUserField(ctx context.Context, db interfaces.DatabaseInterface, user *user_schema.User, field string, value interface{}) error {
	if field == schemas.UserID {
		return fmt.Errorf("UserID cannot be modified")
	}

	if err := user.ShouldHaveFields(schemas.UserID); err != nil {
		return fmt.Errorf("invalid user data: %v", err)
	}

	tag, err := db.Exec(ctx, fmt.Sprintf(`UPDATE users SET %s = $1 WHERE id = $2`, field), value, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user field: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DELETE

func DeleteUser(ctx context.Context, db interfaces.DatabaseInterface, user *user_schema.User) error {
	if err := user.ShouldHaveFields(schemas.UserID); err != nil {
		return fmt.Errorf("invalid user data: %v", err)
	}

	_, err := db.Exec(ctx, `DELETE FROM users WHERE id = $1`, user.ID)
	if err != nil {
		return fmt.Errorf("error deleting user entry: %v", err)
	}

	return nil
}
