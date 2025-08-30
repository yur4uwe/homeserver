package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbselect"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *UserRepo) Query(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	identifier, args := user.Identifier(1)
	if identifier == "" {
		return nil, fmt.Errorf("invalid user data")
	}

	query := qbselect.New().
		Columns(uschema.FieldsList(uconsts.Password)...).
		From(dbconsts.DBUsersTable.Name).
		Where(identifier).
		Build()

	result, err := r.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying user: %v", err)
	}

	return &result, nil
}

func (r *UserRepo) QueryAll(ctx context.Context) ([]uschema.User, error) {
	query := qbselect.New().
		Columns(uschema.FieldsList()...).
		From(dbconsts.DBUsersTable.Name).
		Build()

	results, err := r.db.QueryMany(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}

	return results, nil
}
