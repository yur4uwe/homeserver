package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/constants/dbconsts"
	"homeserver/constants/tconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbselect"
	"homeserver/db/query_builder/qbupdate"
)

func (r *TokenRepo) VerifyToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf(tconsts.ErrEmptyToken)
	}

	query := qbselect.New().
		Columns(uconsts.ID, uconsts.Username).
		From(dbconsts.DBUsersTable.Name).
		Where("auth_token = $1 AND email_verified = false").
		Build()

	user, err := r.db.QueryOne(ctx, query, token)
	if err != nil {
		return fmt.Errorf("error verifying token: %v", err)
	}

	upd_query := qbupdate.New().
		From(dbconsts.DBUsersTable.Name).
		Set("email_verified").
		Where("id = $2").
		Build()

	tag, err := r.db.Pool.Exec(ctx, upd_query, true, user.ID)
	if err != nil {
		return fmt.Errorf("error updating email verification status: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("affected rows mismatch: expected 1, got %d", tag.RowsAffected())
	}

	storagePath := Getenv(constants.StoragePath)

	if err := Mkdir(storagePath+user.Username, 0755); err != nil {
		return fmt.Errorf("error creating user directory: %v", err)
	}

	return nil
}
