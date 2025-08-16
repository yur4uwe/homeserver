package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants/tconsts"
)

func (r *TokenRepo) DeleteToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf(tconsts.ErrEmptyToken)
	}

	tag, err := r.db.Pool.Exec(ctx, `UPDATE users SET auth_token = NULL WHERE auth_token = $1`, token)
	if err != nil {
		return fmt.Errorf("error deleting token: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no token found with value %s", token)
	}

	return nil
}
