package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/db/query_builder/qbselect"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *FolderRepo) GetFolder(ctx context.Context, folder *folschema.UserFolder) (*folschema.UserFolder, error) {
	if folder.ID != 0 {
		return r.GetFolderByID(ctx, folder.ID)
	}

	query := qbselect.New().
		Columns(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("path = $1 AND name = $2").
		Build()

	var folders []folschema.UserFolder
	rows, err := r.db.Pool.Query(ctx, query, folder.Path, folder.Name)
	if err != nil {
		return nil, err
	}

	if err = helpers.ScanRows(rows, &folders); err != nil {
		return nil, err
	}

	if len(folders) != 1 {
		return nil, fmt.Errorf("unexpected error: %d rows returned", len(folders))
	}

	return &folders[0], nil
}

func (r *FolderRepo) GetFolderByID(ctx context.Context, id int) (*folschema.UserFolder, error) {
	query := qbselect.New().
		Columns(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("id = $1").
		Build()

	var folder []folschema.UserFolder
	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if err = helpers.ScanRows(rows, &folder); err != nil {
		return nil, err
	}

	if len(folder) != 1 {
		return nil, fmt.Errorf("unexpected error: %d rows returned", len(folder))
	}

	return &folder[0], nil
}

func (r *FolderRepo) GetFoldersByOwner(ctx context.Context, ownerID int) ([]folschema.UserFolder, error) {
	query := qbselect.New().
		Columns(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("owner_id = $1").
		Build()

	var folders []folschema.UserFolder
	rows, err := r.db.Pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}

	if err = helpers.ScanRows(rows, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *FolderRepo) GetChildren(ctx context.Context, path, name string, owner *uschema.User) ([]folschema.UserFolder, error) {
	parent_query := qbselect.New().
		Columns(folconsts.ID).
		From(dbconsts.DBFoldersTable.Name).
		Where("path = $2 AND name = $3").
		Build()

	query := qbselect.New().
		Columns(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("owner_id = $1 AND parent = (" + parent_query + ")").
		Build()

	rows, err := r.db.Pool.Query(ctx, query, owner.ID, path, name)
	if err != nil {
		return nil, err
	}

	var children []folschema.UserFolder

	if err = helpers.ScanRows(rows, &children); err != nil {
		return nil, err
	}

	return children, nil
}
