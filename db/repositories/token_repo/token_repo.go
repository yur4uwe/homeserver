package tokenrepo

import (
	"context"
	"homeserver/db"
	uschema "homeserver/db/repositories/user_repo/schema"
)

type InTokenRepo interface {
	QueryToken(ctx context.Context, token string) (*uschema.User, error)
}

type TokenRepo struct {
	db *db.Database
}

func NewTokenRepo(db *db.Database) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}
