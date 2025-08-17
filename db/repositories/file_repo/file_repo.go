package filerepo

import (
	"context"
	"homeserver/db"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

type InFileRepo interface {
	DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) error
	UpdateFileEntry(ctx context.Context, modified_file *fschema.UserFile, fields_to_update []string, owner *uschema.User) error
	CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)
}

type FileRepo struct {
	db *db.Database
}

func NewFileRepo(database *db.Database) *FileRepo {
	return &FileRepo{db: database}
}
