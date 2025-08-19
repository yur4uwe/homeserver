package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/folconsts"
	querybuilder "homeserver/db/query_builder"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
	"slices"
	"strconv"
)

func (r *FolderRepo) UpdateFolder(ctx context.Context, folder *folschema.UserFolder, fields_to_update []string, owner *uschema.User) error {
	if slices.Contains(fields_to_update, folconsts.ID) {
		return fmt.Errorf("cannot update folder ID")
	}
	if slices.Contains(fields_to_update, folconsts.OwnerID) {
		return fmt.Errorf("cannot update owner ID")
	}
	if slices.Contains(fields_to_update, folconsts.CreatedAt) {
		return fmt.Errorf("cannot update created at")
	}

	if err := folder.ShouldHaveFields(owner.Username)(fields_to_update...); err != nil {
		return fmt.Errorf("invalid folder data: %v", err)
	}

	fields_str := ""
	for i, field := range fields_to_update {
		if i > 0 {
			fields_str += ", "
		}
		fields_str += fmt.Sprintf("%s = $%d", field, i+1)
	}

	args := make([]any, len(fields_to_update)+1)
	copy(args, OrderUpdateFields(fields_to_update, folder))
	args[len(fields_to_update)] = folder.ID

	query := querybuilder.New().
		Update(dbconsts.DBFoldersTable.Name).
		Set(fields_to_update...).
		Where("id = $" + strconv.Itoa(len(fields_to_update)+1)).
		Returning(folschema.FieldsList()...).
		Build()

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating folder: %v", err)
	}
	defer rows.Close()

	var updated_folder []folschema.UserFolder
	if err = helpers.ScanRows(rows, &updated_folder); err != nil {
		return fmt.Errorf("error scanning folder: %v", err)
	}

	if len(updated_folder) != 1 {
		return fmt.Errorf("error updating folder: %v", err)
	}

	*folder = updated_folder[0]

	return nil
}
