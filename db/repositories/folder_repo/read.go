package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	querybuilder "homeserver/db/query_builder"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/db/repositories/helpers"
)

func (r *FolderRepo) GetFolderByID(ctx context.Context, id int) (*folschema.UserFolder, error) {
	query := querybuilder.New().
		Select(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("id = $1").
		Build()

	var folder []folschema.UserFolder
	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = helpers.ScanRows(rows, &folder); err != nil {
		return nil, err
	}

	if len(folder) != 1 {
		return nil, fmt.Errorf("unexpected error: %d rows returned", len(folder))
	}

	return &folder[0], nil
}

func (r *FolderRepo) GetFolderByPath(ctx context.Context, file *folschema.UserFolder) (*folschema.UserFolder, error) {
	query := querybuilder.New().
		Select(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("path = $1 AND name = $2").
		Build()

	var folder []folschema.UserFolder
	rows, err := r.db.Pool.Query(ctx, query, file.Path, file.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = helpers.ScanRows(rows, &folder); err != nil {
		return nil, err
	}

	if len(folder) != 1 {
		return nil, fmt.Errorf("unexpected error: %d rows returned", len(folder))
	}

	return &folder[0], nil
}

func (r *FolderRepo) GetFoldersByOwner(ctx context.Context, ownerID int) ([]folschema.UserFolder, error) {
	query := querybuilder.New().
		Select(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("owner_id = $1").
		Build()

	var folders []folschema.UserFolder
	rows, err := r.db.Pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = helpers.ScanRows(rows, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *FolderRepo) GetChildren(ctx context.Context, parentID int) ([]folschema.UserFolder, error) {
	query := querybuilder.New().
		Select(folschema.FieldsList()...).
		From(dbconsts.DBFoldersTable.Name).
		Where("parent = $1").
		Build()

	rows, err := r.db.Pool.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []folschema.UserFolder

	if err = helpers.ScanRows(rows, &children); err != nil {
		return nil, err
	}

	return children, nil
}
