package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/constants/dbconsts"
	"homeserver/constants/tconsts"
	"homeserver/constants/uconsts"
	querybuilder "homeserver/db/query_builder"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *TokenRepo) VerifyToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf(tconsts.ErrEmptyToken)
	}

	query := querybuilder.New().
		Select(uconsts.ID, uconsts.Username).
		From("users").
		Where("auth_token = $1 AND email_verified = false").
		Build()

	var user []uschema.User
	// Can fail if the token does not exist == user does not exist
	// Can fail if email is already verified
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

	// Can fail if the username is empty
	if user[0].ID == 0 {
		return fmt.Errorf("invalid user ID for token verification")
	}

	query = querybuilder.New().
		Update(dbconsts.DBUsersTable.Name).
		Set(uconsts.EmailVerified).
		Where("id = $2").
		Build()

	tag, err := r.db.Pool.Exec(ctx, query, true, user[0].ID)
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
