package userrepo

import (
	"homeserver/db"
)

type UserRepo struct {
	db *db.Database
}

func NewUserRepo(db *db.Database) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
