package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/uconsts"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strings"
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

	query := fmt.Sprintf(
		`SELECT %s FROM users WHERE %s`,
		strings.Join(uschema.ReturnableFieldsList(), ", "),
		whereClause,
	)

	var result uschema.User
	err := r.db.Pool.QueryRow(ctx, query, vals...).
		Scan(
			&result.ID,
			&result.Username,
			&result.Email,
			&result.AuthToken,
			&result.CreatedAt,
			&result.EmailVerified,
		)
	if err != nil {
		return nil, fmt.Errorf("error querying user: %v", err)
	}

	return &result, nil
}
