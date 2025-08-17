package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/uconsts"
	uschema "homeserver/db/repositories/user_repo/schema"
	"slices"
	"strings"
)

func (r *UserRepo) Update(ctx context.Context, user *uschema.User, field_vals map[string]any) error {
	fields := make([]string, 0, len(field_vals))
	vals := make([]any, 0, len(field_vals)+1)

	for field := range field_vals {
		if !uschema.IsUserField(field) {
			return fmt.Errorf("invalid field: %s", field)
		}
		if !slices.Contains(uschema.UpdateableFieldsList(), field) {
			return fmt.Errorf("this field is not updatable: %s", field)
		}
		fields = append(fields, field)
		vals = append(vals, field_vals[field])
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := user.ShouldHaveFields(append(fields, uconsts.ID)...); err != nil {
		return fmt.Errorf("invalid user data: %v", err)
	}

	placeholders := ""
	for i := 0; i < len(fields); i++ {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("%s = $%d", fields[i], i+1)
	}

	query := fmt.Sprintf(
		`UPDATE users SET %s WHERE id = $%d RETURNING %s`,
		placeholders, len(fields)+1, strings.Join(uschema.ReturnableFieldsList(), ", "),
	)

	vals = append(vals, user.ID)

	err := r.db.Pool.QueryRow(ctx, query, vals...).
		Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.AuthToken,
			&user.CreatedAt,
			&user.EmailVerified,
		)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}
