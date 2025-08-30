package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/db/query_builder/qbupdate"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *FolderRepo) UpdateFolder(ctx context.Context, folder *folschema.UserFolder, fvp *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	if err := fvp.ShouldNotHaveFields(owner.Username)(folconsts.ID, folconsts.OwnerID, folconsts.CreatedAt); err != nil {
		return nil, fmt.Errorf("inappropriate fields found: %v", err)
	}

	identifier, args := folder.Identificator(owner.Username, 1)
	if identifier == "" {
		return nil, fmt.Errorf("invalid folder data")
	}

	folder.MergeFolders(fvp)

	query := qbupdate.New().
		From(dbconsts.DBFoldersTable.Name).
		Set(folder.GetValidFieldsList()...).
		Where(identifier).
		Returning(folschema.FieldsList()...).
		Build()

	row, err := r.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating folder: %v", err)
	}

	return &row, nil
}
