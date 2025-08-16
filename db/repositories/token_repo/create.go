package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants/tconsts"
)

func (r *TokenRepo) CreateToken(ctx context.Context, userID int, token string) error {
	if userID <= 0 {
		return fmt.Errorf(tconsts.ErrInvalidUserID)
	}
	if token == "" {
		return fmt.Errorf(tconsts.ErrEmptyToken)
	}

	tag, err := r.db.Pool.Exec(ctx, `UPDATE users SET auth_token = $1 WHERE id = $2`, token, userID)
	if err != nil {
		return fmt.Errorf("error creating token: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}
