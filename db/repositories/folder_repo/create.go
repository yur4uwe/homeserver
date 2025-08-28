package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/db/query_builder/qbinsert"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *FolderRepo) CreateFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	if err := folder.ShouldHaveFields(owner.Username)(folconsts.FolderName, folconsts.OwnerID); err != nil {
		return nil, err
	}

	if folder.Parent != nil && *folder.Parent <= 0 {
		return nil, fmt.Errorf("invalid parent folder ID")
	}

	if (folder.Path == "" && folder.Parent != nil) || (folder.Path != "" && folder.Parent == nil) {
		return nil, fmt.Errorf("folder path and parent folder must be set together for non-root folders")
	}

	if folder.Parent == nil && folder.Name != owner.Username {
		return nil, fmt.Errorf("root folder must have the same name as the owner")
	}

	query := qbinsert.New().
		Columns(folschema.FieldsList(folconsts.ID, folconsts.CreatedAt)...).
		Into(dbconsts.DBFoldersTable.Name).
		Returning(folschema.FieldsList()...).
		Build()
	// query := querybuilder.New().
	// 	Insert(folschema.FieldsList(folconsts.ID, folconsts.CreatedAt)...).
	// 	Into(dbconsts.DBFoldersTable.Name).
	// 	Values(len(folschema.FieldsList(folconsts.ID, folconsts.CreatedAt))).
	// 	Returning(folschema.FieldsList()...).
	// 	Build()

	inserted_folder := []folschema.UserFolder{}
	rows, err := r.db.Pool.Query(
		ctx,
		query,
		folder.Name, folder.Path, folder.Parent, folder.OwnerID,
	)
	if err != nil {
		return nil, err
	}

	if err = helpers.ScanRows(rows, &inserted_folder); err != nil {
		return nil, err
	}

	if len(inserted_folder) != 1 {
		return nil, fmt.Errorf("unexpected error: %d rows returned", len(inserted_folder))
	}

	return &inserted_folder[0], nil
}
