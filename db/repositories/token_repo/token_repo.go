package tokenrepo

import (
	"homeserver/db"
)

type InTokenRepo interface {
}

type TokenRepo struct {
	db *db.Database
}

func NewTokenRepo(db *db.Database) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}
