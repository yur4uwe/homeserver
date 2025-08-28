package filerepo

import (
	"context"
	"homeserver/db"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InFileRepo interface {
	CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)

	GetFileByID(ctx context.Context, id int) (*fschema.UserFile, error)
	GetFullFileInfo(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)
	GetFilesByUser(ctx context.Context, owner *uschema.User) ([]fschema.UserFile, error)
	FuzzySearch(ctx context.Context, search_query string, directory string, owner *uschema.User) ([]fschema.UserFile, error)
	GetFilesByDirectory(ctx context.Context, directory string, owner *uschema.User) ([]fschema.UserFile, error)

	DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)
	DeleteFileByID(ctx context.Context, fileID int, owner *uschema.User) (*fschema.UserFile, error)

	UpdateFileEntry(ctx context.Context, file *fschema.UserFile, fvp *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error)
}

type FileRepo struct {
	db *db.Database[fschema.UserFile]
}

func NewFileRepo(pool *pgxpool.Pool) *FileRepo {
	database := db.Database[fschema.UserFile]{Pool: pool}
	return &FileRepo{db: &database}
}
