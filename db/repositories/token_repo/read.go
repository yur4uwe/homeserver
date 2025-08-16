package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants/tconsts"
	uschema "homeserver/db/repositories/user_repo/schema"
	"os"
)

var (
	Mkdir  = os.Mkdir
	Getenv = os.Getenv
)

func (r *TokenRepo) CheckToken(ctx context.Context, token string) (int, error) {
	if token == "" {
		return 0, fmt.Errorf(tconsts.ErrEmptyToken)
	}

	var userID int
	err := r.db.Pool.QueryRow(ctx, `SELECT id FROM users WHERE auth_token = $1 AND email_verified = true`, token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error checking token: %v", err)
	}
	return userID, nil
}

func (r *TokenRepo) QueryToken(ctx context.Context, token string) (*uschema.User, error) {
	if token == "" {
		return nil, fmt.Errorf(tconsts.ErrEmptyToken)
	}

	var user uschema.User
	err := r.db.Pool.QueryRow(ctx, `
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
