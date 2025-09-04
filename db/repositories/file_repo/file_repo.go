package filerepo

import (
	"context"
	"homeserver/db"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/interfaces"
)

type InFileRepo interface {
	CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)

	GetFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)
	GetFilesByUser(ctx context.Context, owner *uschema.User) ([]fschema.UserFile, error)
	FuzzySearch(ctx context.Context, search_query string, directory string, owner *uschema.User) ([]fschema.UserFile, error)
	GetFilesByDirectory(ctx context.Context, directory string, owner *uschema.User) ([]fschema.UserFile, error)

	DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)

	UpdateFileEntry(ctx context.Context, file *fschema.UserFile, fvp *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)
}

type FileRepo struct {
	db *db.Database[fschema.UserFile]
}

func NewFileRepo(pool interfaces.QueryExecutor) *FileRepo {
	database := db.Database[fschema.UserFile]{Pool: pool}
	return &FileRepo{db: &database}
}
