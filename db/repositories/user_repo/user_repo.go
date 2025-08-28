package userrepo

import (
	"context"
	"homeserver/db"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InUserRepo interface {
	Create(ctx context.Context, user *uschema.User) (*uschema.User, error)
	Query(ctx context.Context, user *uschema.User) (*uschema.User, error)
	QueryAll(ctx context.Context) ([]uschema.User, error)
	Update(ctx context.Context, user *uschema.User, field_vals map[string]any) error
	Delete(ctx context.Context, user *uschema.User) error
}

type UserRepo struct {
	db *db.Database[uschema.User]
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	database := db.Database[uschema.User]{Pool: pool}

	return &UserRepo{
		db: &database,
	}
}
