package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	querybuilder "homeserver/db/query_builder"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *UserRepo) Query(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	validFields := user.GetValidFields()
	if len(validFields) == 1 && validFields[0] == uconsts.EmailVerified {
		return nil, fmt.Errorf("no valid fields to query")
	}

	whereClause := ""
	vals := make([]any, 0, len(validFields))
	for i, field := range validFields {
		if i > 0 {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%s = $%d", field, i+1)
		if value, ok := user.Get(field); ok {
			vals = append(vals, value)
		}
	}

	query := querybuilder.New().
		Select(uschema.FieldsList(uconsts.Password)...).
		From(dbconsts.DBUsersTable.Name).
		Where(whereClause).
		Build()

	// query := fmt.Sprintf(
	// 	`SELECT %s FROM users WHERE %s`,
	// 	strings.Join(uschema.ReturnableFieldsList(), ", "),
	// 	whereClause,
	// )

	var result []uschema.User
	rows, err := r.db.Pool.Query(ctx, query, vals...)
	if err != nil {
		return nil, fmt.Errorf("error querying user: %v", err)
	}
	defer rows.Close()

	if err = helpers.ScanRows(rows, &result); err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("unexpected error: %d rows returned", len(result))
	}

	return &result[0], nil
}
