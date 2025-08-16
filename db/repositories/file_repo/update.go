package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/fconsts"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"slices"
	"strings"
)

var (
	FixDuplicates = PreventDuplicateFileName
)

func (repo *FileRepo) UpdateFileEntry(ctx context.Context, modified_file *fschema.UserFile, fields_to_update []string, owner *uschema.User) error {
	if slices.Contains(fields_to_update, fconsts.ID) {
		return fmt.Errorf("cannot update file ID")
	}
	if slices.Contains(fields_to_update, fconsts.OwnerID) {
		return fmt.Errorf("cannot update owner ID")
	}
	if slices.Contains(fields_to_update, fconsts.CreatedAt) {
		return fmt.Errorf("cannot update created at")
	}
	if slices.Contains(fields_to_update, fconsts.UpdatedAt) {
		return fmt.Errorf("cannot update updated at")
	}

	args := append(fields_to_update, fconsts.ID)
	if err := modified_file.ShouldHaveFields(owner.Username)(args...); err != nil {
		return fmt.Errorf("invalid file data: %v", err)
	}

	var err error
	if modified_file.Filename != "" && modified_file.Extension != "" && modified_file.Path != "" {
		modified_file, err = FixDuplicates(modified_file)
		if err != nil {
			return fmt.Errorf("error preventing duplicate file name: %v", err)
		}
	} else if modified_file.Filename != "" || modified_file.Extension != "" || modified_file.Path != "" {
		return fmt.Errorf("if filename, extension, or path is provided for file update, all must be provided for the update")
	}

	fields_str := ""
	for i, field := range fields_to_update {
		if i > 0 {
			fields_str += ", "
		}
		fields_str += fmt.Sprintf("%s = $%d", field, i+1)
	}

	query := fmt.Sprintf(`
		UPDATE user_files
		SET %s
		WHERE id = $%d
		RETURNING`+strings.Join(fschema.FieldsList(), ", "),
		fields_str, len(fields_to_update)+1)

	query_args := make([]any, len(fields_to_update)+1)
	copy(query_args, OrderFileUpdateFields(fields_to_update, modified_file))
	query_args[len(fields_to_update)] = modified_file.ID

	err = repo.db.Pool.QueryRow(
		ctx,
		query,
		query_args...,
	).Scan(
		&modified_file.ID, &modified_file.Filename, &modified_file.Extension,
		&modified_file.Path, &modified_file.Directory, &modified_file.Size,
		&modified_file.ContentType, &modified_file.OwnerID, &modified_file.CreatedAt,
		&modified_file.UpdatedAt, &modified_file.Metadata,
	)
	if err != nil {
		return fmt.Errorf("error updating file entry: %v", err)
	}

	return nil
}
