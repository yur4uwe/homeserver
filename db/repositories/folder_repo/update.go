package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/db/query_builder/qbupdate"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strconv"
)

func (r *FolderRepo) UpdateFolder(ctx context.Context, folder *folschema.UserFolder, fvp *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	if err := fvp.ShouldNotHaveFields(owner.Username)(folconsts.ID, folconsts.OwnerID, folconsts.CreatedAt); err != nil {
		return nil, fmt.Errorf("inappropriate fields found: %v", err)
	}

	if err := folder.ShouldHaveFields(owner.Username)(folconsts.ID); err != nil {
		return nil, fmt.Errorf("invalid folder data: %v", err)
	}

	folder.MergeFolders(fvp)

	args := make([]any, 0, len(folder.GetValidFieldsList())+1)
	args = append(args, folder.GetValidFieldsData()...)
	args = append(args, folder.ID)

	query := qbupdate.New().
		From(dbconsts.DBFoldersTable.Name).
		Set(folder.GetValidFieldsList()...).
		Where("id = $" + strconv.Itoa(len(folder.GetValidFieldsList())+1)).
		Returning(folschema.FieldsList()...).
		Build()

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating folder: %v", err)
	}

	var updated_folder []folschema.UserFolder
	if err = helpers.ScanRows(rows, &updated_folder); err != nil {
		return nil, fmt.Errorf("error scanning folder: %v", err)
	}

	if len(updated_folder) != 1 {
		return nil, fmt.Errorf("error updating folder: %v", err)
	}

	return &updated_folder[0], nil
}
