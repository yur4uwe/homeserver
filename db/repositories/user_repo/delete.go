package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants/uconsts"
	querybuilder "homeserver/db/query_builder"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *UserRepo) Delete(ctx context.Context, user *uschema.User) error {
	if err := user.ShouldHaveFields(uconsts.ID); err != nil {
		return fmt.Errorf("invalid user data: %v", err)
	}

	query := querybuilder.New().Delete().From("users").Where("id = $1").Build()

	tag, err := r.db.Pool.Exec(ctx, query, user.ID)
	if err != nil {
		return fmt.Errorf("error deleting user entry: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
