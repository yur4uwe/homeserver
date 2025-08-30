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
	user_fields := user.GetValidFields()

	whereClause := ""
	for i, field := range user_fields {
		if i > 0 {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%s = $%d", field, i+1)
	}

	query := qbselect.New().
		Columns(uschema.FieldsList(uconsts.Password)...).
		From(dbconsts.DBUsersTable.Name).
		Where(whereClause).
		Build()

	result, err := r.db.QueryOne(ctx, query, user.GetValidFieldsData()...)
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
