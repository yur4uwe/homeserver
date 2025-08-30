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

	valid_fields := fvp.GetValidFieldsList()

	identifier, args := folder.Identificator(owner.Username, len(valid_fields)+1)
	if identifier == "" {
		return nil, fmt.Errorf("invalid folder data")
	}
	folder.MergeFolders(fvp)

	query := qbupdate.New().
		From(dbconsts.DBFoldersTable.Name).
		Set(valid_fields...).
		Where(identifier).
		Returning(folschema.FieldsList()...).
		Build()

	qargs := make([]any, 0, len(valid_fields)+len(args))
	qargs = append(qargs, fvp.GetValidFieldsData()...)
	qargs = append(qargs, args...)

	row, err := r.db.QueryOne(ctx, query, qargs...)
	if err != nil {
		return nil, fmt.Errorf("error updating folder: %v", err)
	}

	return &row, nil
}
