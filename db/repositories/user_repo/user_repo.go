package userrepo

import (
	"homeserver/db"
)

type InUserRepo interface {
}

type UserRepo struct {
	db *db.Database
}

func NewUserRepo(db *db.Database) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
