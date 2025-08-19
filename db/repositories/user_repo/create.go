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

func (r *UserRepo) Create(ctx context.Context, user *uschema.User) (*uschema.User, error) {
	if err := user.ShouldHaveFields(uconsts.Username, uconsts.Password, uconsts.Email); err != nil {
		return nil, fmt.Errorf("invalid user data: %v", err)
	}

	query := querybuilder.New().
		Insert(uschema.FieldsList(uconsts.ID, uconsts.CreatedAt)...).
		Into(dbconsts.DBUsersTable.Name).
		Values(len(uschema.FieldsList(uconsts.ID, uconsts.CreatedAt))).
		Returning(uschema.FieldsList(uconsts.Password)...).
		Build()

	// query := fmt.Sprintf(
	// 	`INSERT INTO users (%s, created_at) VALUES (%s, now()) RETURNING %s`,
	// 	strings.Join(uschema.FieldsList(uconsts.ID, uconsts.UpdatedAt, uconsts.CreatedAt), ", "),
	// 	placeholders,
	// 	strings.Join(uschema.FieldsList(uconsts.Password), ", "))
	rows, err := r.db.Pool.Query(ctx, query, user.Username, user.Password, user.Email, nil, false)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}
	defer rows.Close()

	var inserted_user []uschema.User
	if err = helpers.ScanRows(rows, &inserted_user); err != nil {
		return nil, err
	}

	if len(inserted_user) > 1 {
		return nil, fmt.Errorf("unexpected error: multiple rows returned")
	}

	return &inserted_user[0], nil
}
