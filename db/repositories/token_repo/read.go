package tokenrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/tconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbselect"
	uschema "homeserver/db/repositories/user_repo/schema"
	"os"
)

var (
	Mkdir  = os.Mkdir
	Getenv = os.Getenv
)

func (r *TokenRepo) QueryToken(ctx context.Context, token string, should_be_verified bool) (*uschema.User, error) {
	if token == "" {
		return nil, fmt.Errorf(tconsts.ErrEmptyToken)
	}

	query := qbselect.New().
		Columns(uschema.FieldsList(uconsts.Password)...).
		From(dbconsts.DBUsersTable.Name).
		Where("auth_token = $1 AND email_verified = $2").
		Build()

	user, err := r.db.QueryOne(ctx, query, token, should_be_verified)
	if err != nil {
		return nil, fmt.Errorf("error querying token: %v", err)
	}

	return &user, nil
}
