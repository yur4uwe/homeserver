package filerepo

import (
	"homeserver/db"
)

type InFileRepo interface {
}

type FileRepo struct {
	db *db.Database
}

func NewFileRepo(database *db.Database) *FileRepo {
	return &FileRepo{db: database}
}
