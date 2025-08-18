package folderrepo

import (
	"context"
	"homeserver/constants/folconsts"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"strings"
)

func (r *FolderRepo) GetFolderByID(ctx context.Context, id int) (*folschema.UserFolder, error) {
	fields_to_read := strings.Join(folschema.FieldsList(folconsts.CreatedAt), ", ")

	var folder folschema.UserFolder
	err := r.db.Pool.QueryRow(ctx, "SELECT "+fields_to_read+" FROM folders WHERE id = ?", id).
		Scan(&folder.ID, &folder.Name, &folder.Path, &folder.Parent, &folder.OwnerID)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepo) GetFolderByPath(ctx context.Context, file *folschema.UserFolder) (*folschema.UserFolder, error) {
	fields_to_read := strings.Join(folschema.FieldsList(folconsts.CreatedAt), ", ")

	var folder folschema.UserFolder
	err := r.db.Pool.QueryRow(ctx, "SELECT "+fields_to_read+" FROM folders WHERE path = ? AND name = ?", file.Path, file.Name).
		Scan(&folder.ID, &folder.Name, &folder.Path, &folder.Parent, &folder.OwnerID)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepo) GetFolderByOwner(ctx context.Context, ownerID int) (*folschema.UserFolder, error) {
	fields_to_read := strings.Join(folschema.FieldsList(folconsts.CreatedAt), ", ")

	var folder folschema.UserFolder
	err := r.db.Pool.QueryRow(ctx, "SELECT "+fields_to_read+" FROM folders WHERE owner_id = ?", ownerID).
		Scan(&folder.ID, &folder.Name, &folder.Path, &folder.Parent, &folder.OwnerID)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *FolderRepo) GetChildren(ctx context.Context, parentID int) ([]*folschema.UserFolder, error) {
	fields_to_read := strings.Join(folschema.FieldsList(folconsts.CreatedAt), ", ")

	rows, err := r.db.Pool.Query(ctx, "SELECT "+fields_to_read+" FROM folders WHERE parent = ?", parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []*folschema.UserFolder
	for rows.Next() {
		var folder folschema.UserFolder
		if err := rows.Scan(&folder.ID, &folder.Name, &folder.Path, &folder.Parent, &folder.OwnerID); err != nil {
			return nil, err
		}
		children = append(children, &folder)
	}
	return children, nil
}
