package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	"homeserver/db/query_builder/qbdelete"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *FolderRepo) DeleteFolderByID(ctx context.Context, folderID int, ownerID int) (*folschema.UserFolder, error) {
	if folderID <= 0 {
		return nil, fmt.Errorf("invalid folder ID")
	}

	query := qbdelete.New().From(dbconsts.DBUsersTable.Name).Where("id = $1 AND owner_id = $2").Returning(folschema.FieldsList()...).Build()

	fol, err := r.db.QueryOne(ctx, query, folderID, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error deleting folder entry: %v", err)
	}

	return &fol, nil
}

func (r *FolderRepo) DeleteFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	if folder.ShouldHaveFields(owner.Username)(folconsts.Path, folconsts.FolderName, folconsts.OwnerID) != nil {
		return nil, fmt.Errorf("invalid folder data")
	}

	valid_fields := folder.GetValidFieldsList()

	whereClause := "WHERE "
	for i, field := range valid_fields {
		if i > 0 {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%s = $%d", field, i+1)
	}

	args := make([]any, 0, len(valid_fields))
	args = append(args, folder.GetValidFieldsData()...)

	query := qbdelete.New().From(dbconsts.DBFoldersTable.Name).Where(whereClause).Returning(folschema.FieldsList()...).Build()

	fol, err := r.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error deleting folder entry: %v", err)
	}

	return &fol, nil
}
