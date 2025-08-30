package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbinsert"
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

	inserted_user, err := r.db.QueryOne(ctx, query, user.Username, user.Password, user.Email, nil, false)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return &inserted_user, nil
}
