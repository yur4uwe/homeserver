package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/constants/tconsts"
)

func (r *TokenRepo) VerifyToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf(tconsts.ErrEmptyToken)
	}

	var id int
	var userName string
	// Can fail if the token does not exist == user does not exist
	// Can fail if email is already verified
	err := r.db.Pool.QueryRow(ctx, `SELECT id, username FROM users WHERE auth_token = $1 AND email_verified = false`, token).Scan(&id, &userName)
	if err != nil {
		return fmt.Errorf("error verifying token: %v", err)
	}

	// Can fail if the username is empty
	if id == 0 {
		return fmt.Errorf("invalid user ID for token verification")
	}

	tag, err := r.db.Pool.Exec(ctx, `UPDATE users SET email_verified = true WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error updating email verification status: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("affected rows mismatch: expected 1, got %d", tag.RowsAffected())
	}

	storagePath := Getenv(constants.StoragePath)

	if err := Mkdir(storagePath+userName, 0755); err != nil {
		return fmt.Errorf("error creating user directory: %v", err)
	}

	return nil
}
