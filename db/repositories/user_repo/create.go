package userrepo

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/constants/uconsts"
	querybuilder "homeserver/db/query_builder"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *UserRepo) Create(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	if err := user.ShouldHaveFields(uconsts.Username, uconsts.Password, uconsts.Email); err != nil {
		return nil, fmt.Errorf("invalid user data: %v", err)
	}

	var inserted_user uschema.User
	query := querybuilder.New().
		Insert(uschema.FieldsList(uconsts.ID)...).
		Into(constants.DBUsersTable.Name).
		Values(user.Username, user.Email, user.Password, user.CreatedAt).
		Returning(uschema.FieldsList(uconsts.Password)...).
		Build()

	// query := fmt.Sprintf(
	// 	`INSERT INTO users (%s, created_at) VALUES (%s, now()) RETURNING %s`,
	// 	strings.Join(uschema.FieldsList(uconsts.ID, uconsts.UpdatedAt, uconsts.CreatedAt), ", "),
	// 	placeholders,
	// 	strings.Join(uschema.FieldsList(uconsts.Password), ", "))
	if err := r.db.Pool.QueryRow(ctx, query, user.Username, user.Password, user.Email, nil, false).
		Scan(&inserted_user.ID, &inserted_user.Username, &inserted_user.Email); err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return &inserted_user, nil
}
