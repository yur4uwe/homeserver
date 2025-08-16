package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/folconsts"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *FolderRepo) DeleteFolderByID(ctx context.Context, folderID int, ownerID int) (int, error) {
	if folderID <= 0 {
		return 0, fmt.Errorf("invalid folder ID")
	}

	var deletedID int
	err := r.db.Pool.QueryRow(
		ctx,
		`DELETE FROM folders WHERE id = $1 AND owner_id = $2 RETURNING id`,
		folderID,
		ownerID,
	).Scan(&deletedID)
	if err != nil {
		return 0, fmt.Errorf("error deleting folder entry: %v", err)
	}

	return deletedID, nil
}

func (r *FolderRepo) DeleteFolderByPath(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (int, error) {
	if folder.ShouldHaveFields(owner.Username)(folconsts.Path, folconsts.FolderName, folconsts.OwnerID) != nil {
		return 0, fmt.Errorf("invalid folder data")
	}

	var deletedID int
	err := r.db.Pool.QueryRow(
		ctx,
		`DELETE FROM folders WHERE path = $1 AND name = $2 AND owner_id = $3 RETURNING id`,
		folder.Path,
		folder.Name,
		owner.ID,
	).Scan(&deletedID)
	if err != nil {
		return 0, fmt.Errorf("error deleting folder entry: %v", err)
	}

	return deletedID, nil
}
