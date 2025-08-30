package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/db/query_builder/qbdelete"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (r *FolderRepo) DeleteFolder(ctx context.Context, folder *folschema.UserFolder, owner *uschema.User) (*folschema.UserFolder, error) {
	identifier, args := folder.Identificator(owner.Username, 1)
	if identifier == "" {
		return nil, fmt.Errorf("invalid folder data")
	}

	query := qbdelete.New().
		From(dbconsts.DBFoldersTable.Name).
		Where(identifier).
		Returning(folschema.FieldsList()...).
		Build()

	fol, err := r.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error deleting folder entry: %v", err)
	}

	return &fol, nil
}
