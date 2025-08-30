package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/db/query_builder/qbselect"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *FolderRepo) GetFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	uname := ""
	if owner != nil {
		uname = owner.Username
	}

	identifier, args := folder.Identificator(uname, 1)
	if identifier == "" {
		return nil, fmt.Errorf("folder is not identifiable")
	}

	query := qbselect.New().
		Columns(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where(identifier).
		Build()

	full_folder, err := r.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &full_folder, nil
}

func (r *FolderRepo) GetFoldersByOwner(ctx context.Context, ownerID int) ([]folschema.UserFolder, error) {
	query := qbselect.New().
		Columns(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("owner_id = $1").
		Build()

	folders, err := r.db.QueryMany(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error querying folders: %v", err)
	}

	return folders, nil
}

func (r *FolderRepo) GetChildren(ctx context.Context, path, name string, owner *uschema.User) ([]folschema.UserFolder, error) {
	parent_query := qbselect.New().
		Columns(folconsts.ID).
		From(dbconsts.DBFoldersTable.Name).
		Where("path = $2 AND name = $3 AND owner_id = $4").
		Build()

	query := qbselect.New().
		Columns(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("owner_id = $1 AND parent = (" + parent_query + ")").
		Build()

	children, err := r.db.QueryMany(ctx, query, owner.ID, path, name, owner.ID)
	if err != nil {
		return nil, err
	}

	return children, nil
}
