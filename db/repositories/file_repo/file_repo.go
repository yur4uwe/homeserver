package filerepo

import (
	"context"
	"homeserver/db"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

type InFileRepo interface {
	CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)

	GetFileByID(ctx context.Context, id string) (*fschema.UserFile, error)
	GetFullFileInfo(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)
	GetFilesByUser(ctx context.Context, owner *uschema.User) ([]fschema.UserFile, error)
	FuzzySearch(ctx context.Context, search_query string, directory string, owner *uschema.User) ([]fschema.UserFile, error)
	GetFilesByDirectory(ctx context.Context, directory string, owner *uschema.User) ([]fschema.UserFile, error)

	DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (int, error)
	DeleteFileByID(ctx context.Context, fileID string, owner *uschema.User) error

	UpdateFileEntry(ctx context.Context, file *fschema.UserFile, fvp map[string]any, owner *uschema.User) error
}

type FileRepo struct {
	db *db.Database
}

func NewFileRepo(database *db.Database) *FileRepo {
	return &FileRepo{db: database}
}
