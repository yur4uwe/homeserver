package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/db/query_builder/qbdelete"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *UserRepo) Delete(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	whereClause, args := user.Identifier(1)
	if whereClause == "" {
		return nil, fmt.Errorf("invalid user data")
	}

	query := qbdelete.New().
		From(dbconsts.DBUsersTable.Name).
		Where(whereClause).
		Returning(uschema.FieldsList()...).
		Build()

	deleted_user, err := r.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error deleting user entry: %v", err)
	}

	return &deleted_user, nil
}
