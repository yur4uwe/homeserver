package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/uconsts"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strings"
)

func (r *UserRepo) Create(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	if err := user.ShouldHaveFields(uconsts.Username, uconsts.Password, uconsts.Email); err != nil {
		return nil, fmt.Errorf("invalid user data: %v", err)
	}

	placeholders := ""
	for i := 0; i < len(uschema.UpdateableFieldsList()); i++ {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
	}

	var inserted_user uschema.User
	query := fmt.Sprintf(
		`INSERT INTO users (%s, created_at) VALUES (%s, now()) RETURNING %s`,
		strings.Join(uschema.UpdateableFieldsList(), ", "), placeholders, strings.Join(uschema.ReturnableFieldsList(), ", "))
	if err := r.db.Pool.QueryRow(ctx, query, user.Username, user.Password, user.Email, nil, false).
		Scan(&inserted_user.ID, &inserted_user.Username, &inserted_user.Email); err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return &inserted_user, nil
}
