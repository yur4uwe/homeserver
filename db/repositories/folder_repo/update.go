package folderrepo

import (
	"context"
	"fmt"
	"homeserver/constants/folconsts"
	folschema "homeserver/db/repositories/folder_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"slices"
	"strings"
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

	err := r.db.Pool.QueryRow(
		ctx,
		fmt.Sprintf(
			`UPDATE folders SET %s WHERE id = $%d RETURNING `+strings.Join(folschema.FieldsList(true), ", "),
			fields_str, len(fields_to_update)+1),
		args...,
	).Scan(
		&folder.ID,
		&folder.Path,
		&folder.Name,
		&folder.Parent,
		&folder.OwnerID,
		&folder.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error updating folder: %v", err)
	}

	return nil
}
