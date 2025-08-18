package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants"
	"homeserver/constants/folconsts"
	querybuilder "homeserver/db/query_builder"
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

	query := querybuilder.New().
		Insert(folschema.FieldsList(folconsts.ID, folconsts.CreatedAt)...).
		Into(constants.DBFoldersTable.Name).
		Values(folder.Name, folder.Path, folder.Parent, folder.OwnerID).
		Build()

	inserted_folder := folschema.UserFolder{}
	rows, err := r.db.Pool.Query(
		ctx,
		query,
		folder.Name, folder.Path, folder.Parent, folder.OwnerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("failed to insert folder: no rows returned")
	}

	err = helpers.ScanRows(
		rows,
		&inserted_folder,
	)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return nil, fmt.Errorf("unexpected error: multiple rows returned")
	}
	return &inserted_folder, nil
}
