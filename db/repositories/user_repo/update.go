package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbupdate"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strings"
)

func (r *UserRepo) Update(ctx context.Context, user *uschema.User, fvp *uschema.User) (*uschema.User, error) {
	if err := fvp.ShouldNotHaveFields(uconsts.ID, uconsts.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to validate fields: %v", err)
	}

	valid_fields := fvp.GetValidFields()

	identifier, args := user.Identifier(len(valid_fields) + 1)
	if !strings.Contains(identifier, uconsts.ID) {
		return nil, fmt.Errorf("invalid user identifier")
	}

	user.Merge(fvp)

	query := qbupdate.New().
		From(dbconsts.DBUsersTable.Name).
		Set(valid_fields...).
		Where(identifier).
		Build()

	qargs := make([]any, 0, len(args)+len(valid_fields))
	qargs = append(qargs, args...)
	qargs = append(qargs, fvp.GetValidFieldsData()...)

	updated_user, err := r.db.QueryOne(ctx, query, qargs...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &updated_user, nil
}
