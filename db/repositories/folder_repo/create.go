package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/db/query_builder/qbinsert"
	folschema "homeserver/db/repositories/folder_repo/schema"
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
		Columns(folschema.FieldsList(folconsts.ID, folconsts.CreatedAt, folconsts.TrustedSource)...).
		Into(dbconsts.DBFoldersTable.Name).
		Returning(folschema.FieldsList()...).
		OnConflict([]string{folconsts.FolderName, folconsts.Path, folconsts.OwnerID}, "DO UPDATE SET name = EXCLUDED.name").
		Build()

	inserted_folder, err := r.db.QueryOne(ctx, query, folder.Name, folder.Path, folder.Parent, folder.OwnerID)
	if err != nil {
		return nil, err
	}

	return &inserted_folder, nil
}
