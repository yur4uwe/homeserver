package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/folconsts"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strings"
)

func (r *FolderRepo) CreateFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	write_fields := []string{
		folconsts.FolderName,
		folconsts.Path,
		folconsts.Parent,
		folconsts.OwnerID,
		folconsts.CreatedAt,
	}

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

	fields_to_read := strings.Join(write_fields, ", ")

	placeholders := strings.Repeat("?, ", folschema.FieldNum(false)-1) + "now()"

	inserted_folder := folschema.UserFolder{}
	err := r.db.Pool.QueryRow(
		ctx,
		"INSERT INTO folders ("+fields_to_read+") VALUES ("+placeholders+")",
		folder.Name, folder.Path, folder.Parent, folder.OwnerID,
	).Scan(
		&inserted_folder.ID,
		&inserted_folder.Name,
		&inserted_folder.Path,
		&inserted_folder.Parent,
		&inserted_folder.OwnerID,
		&inserted_folder.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &inserted_folder, nil
}
