package folderrepo

import (
	"context"
	"homeserver/db"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

type InFolderRepo interface {
	GetFolderByID(ctx context.Context, id int) (*folschema.UserFolder, error)
	GetFolderByPath(ctx context.Context, folder *folschema.UserFolder) (*folschema.UserFolder, error)
	GetFolderByOwner(ctx context.Context, ownerID int) (*folschema.UserFolder, error)
	GetChildren(ctx context.Context, parentID int) ([]*folschema.UserFolder, error)
	CreateFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error)
	UpdateFolder(ctx context.Context, folder *folschema.UserFolder, fieldsToUpdate []string, owner *uschema.User) error
	DeleteFolderByID(ctx context.Context, id int, ownerID int) (int, error)
	DeleteFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (int, error)
}

type FolderRepo struct {
	db *db.Database
}

func NewFolderRepo(database *db.Database) InFolderRepo {
	return &FolderRepo{
		db: database,
	}
}
