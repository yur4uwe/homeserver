package userrepo

import (
	"context"
	"homeserver/db"
	uschema "homeserver/db/repositories/user_repo/schema"
)

type InUserRepo interface {
	Create(ctx context.Context, user *uschema.User) (*uschema.User, error)
	Query(ctx context.Context, user *uschema.User) (*uschema.User, error)
	QueryAll(ctx context.Context) ([]uschema.User, error)
	Update(ctx context.Context, user *uschema.User, field_vals map[string]any) error
	Delete(ctx context.Context, user *uschema.User) error
}

type UserRepo struct {
	db *db.Database
}

func NewUserRepo(db *db.Database) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
