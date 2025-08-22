package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbinsert"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *UserRepo) Create(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	if err := user.ShouldHaveFields(uconsts.Username, uconsts.Password, uconsts.Email); err != nil {
		return nil, fmt.Errorf("invalid user data: %v", err)
	}

	query := qbinsert.New().
		Columns(uschema.FieldsList(uconsts.ID, uconsts.CreatedAt)...).
		Into(dbconsts.DBUsersTable.Name).
		Returning(uschema.FieldsList(uconsts.Password)...).
		Build()

	rows, err := r.db.Pool.Query(ctx, query, user.Username, user.Password, user.Email, nil, false)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}
	defer rows.Close()

	var inserted_user []uschema.User
	if err = helpers.ScanRows(rows, &inserted_user); err != nil {
		return nil, err
	}

	if len(inserted_user) > 1 {
		return nil, fmt.Errorf("unexpected error: multiple rows returned")
	}

	return &inserted_user[0], nil
}
