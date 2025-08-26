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
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
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

	var user []uschema.User
	rows, err := r.db.Pool.Query(ctx, query, token)
	if err != nil {
		return fmt.Errorf("error verifying token: %v", err)
	}
	defer rows.Close()

	if err := helpers.ScanRows(rows, &user); err != nil {
		return fmt.Errorf("error scanning rows: %v", err)
	}

	if len(user) != 1 {
		return fmt.Errorf("unexpected error: %v", err)
	}

	if user[0].ID == 0 {
		return fmt.Errorf("invalid user ID for token verification")
	}

	upd_query := qbupdate.New().
		From(dbconsts.DBUsersTable.Name).
		Set("email_verified").
		Where("id = $2").
		Build()

	tag, err := r.db.Pool.Exec(ctx, upd_query, true, user[0].ID)
	if err != nil {
		return fmt.Errorf("error updating email verification status: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("affected rows mismatch: expected 1, got %d", tag.RowsAffected())
	}

	storagePath := Getenv(constants.StoragePath)

	if err := Mkdir(storagePath+user[0].Username, 0755); err != nil {
		return fmt.Errorf("error creating user directory: %v", err)
	}

	return nil
}
