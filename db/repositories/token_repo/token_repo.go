package tokenrepo

import (
	"context"
	"homeserver/db"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InTokenRepo interface {
	QueryToken(ctx context.Context, token string) (*uschema.User, error)
	CreateToken(ctx context.Context, userID int, token string) error
	DeleteToken(ctx context.Context, token string) error
	CheckToken(ctx context.Context, token string) (int, error)
	VerifyToken(ctx context.Context, token string) error
}

type TokenRepo struct {
	db *db.Database[uschema.User]
}

func NewTokenRepo(pool *pgxpool.Pool) *TokenRepo {
	database := db.Database[uschema.User]{Pool: pool}

	return &TokenRepo{
		db: &database,
	}
}
