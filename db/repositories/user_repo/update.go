package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	querybuilder "homeserver/db/query_builder"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
	"slices"
)

func (r *UserRepo) Update(ctx context.Context, user *uschema.User, field_vals map[string]any) error {
	fields := make([]string, 0, len(field_vals))
	vals := make([]any, 0, len(field_vals)+1)

	for field := range field_vals {
		if !uschema.IsUserField(field) {
			return fmt.Errorf("invalid field: %s", field)
		}
		if !slices.Contains(uschema.FieldsList(uconsts.ID, uconsts.CreatedAt), field) {
			return fmt.Errorf("this field is not updatable: %s", field)
		}
		fields = append(fields, field)
		vals = append(vals, field_vals[field])
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := user.ShouldHaveFields(uconsts.ID); err != nil {
		return fmt.Errorf("invalid user data: %v", err)
	}

	query := querybuilder.New().
		Update(dbconsts.DBUsersTable.Name).
		Set(fields...).
		Where(fmt.Sprintf("id = $%d", len(fields)+1)).
		Returning(uschema.FieldsList(uconsts.Password)...).
		Build()

	vals = append(vals, user.ID)

	rows, err := r.db.Pool.Query(ctx, query, vals...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	defer rows.Close()

	var updated_user []uschema.User
	if err = helpers.ScanRows(rows, &updated_user); err != nil {
		return err
	}

	if len(updated_user) != 1 {
		return fmt.Errorf("unexpected error: %d rows returned", len(updated_user))
	}

	*user = updated_user[0]

	return nil
}
