package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/tconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbselect"
	"homeserver/db/repositories/helpers"
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

	query := qbselect.New().
		Columns(uconsts.ID).
		From(dbconsts.DBUsersTable.Name).
		Where("auth_token = $1 AND email_verified = true").
		Build()

	var userID int
	err := r.db.Pool.QueryRow(ctx, query, token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error checking token: %v", err)
	}
	return userID, nil
}

func (r *TokenRepo) QueryToken(ctx context.Context, token string) (*uschema.User, error) {
	if token == "" {
		return nil, fmt.Errorf(tconsts.ErrEmptyToken)
	}

	query := qbselect.New().
		Columns(uschema.FieldsList(uconsts.Password)...).
		From(dbconsts.DBUsersTable.Name).
		Where("auth_token = $1").
		Build()

	var user []uschema.User
	rows, err := r.db.Pool.Query(ctx, query, token)
	if err != nil {
		return nil, fmt.Errorf("error querying token: %v", err)
	}
	defer rows.Close()

	if err = helpers.ScanRows(rows, &user); err != nil {
		return nil, fmt.Errorf("error querying token: %v", err)
	}

	if len(user) != 1 {
		return nil, fmt.Errorf("unexpected error: %d users found", len(user))
	}

	return &user[0], nil
}
