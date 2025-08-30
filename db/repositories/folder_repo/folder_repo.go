package folderrepo

import (
	"context"
	"homeserver/db"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InFolderRepo interface {
	GetFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error)
	GetFoldersByOwner(ctx context.Context, ownerID int) ([]folschema.UserFolder, error)
	GetChildren(ctx context.Context, path, name string, owner *uschema.User) ([]folschema.UserFolder, error)
	CreateFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error)
	UpdateFolder(ctx context.Context, folder *folschema.UserFolder, fvp *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error)
	DeleteFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error)
}

type FolderRepo struct {
	db *db.Database[folschema.UserFolder]
}

func NewFolderRepo(pool *pgxpool.Pool) InFolderRepo {
	database := db.Database[folschema.UserFolder]{Pool: pool}

	return &FolderRepo{
		db: &database,
	}
}
