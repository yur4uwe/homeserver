package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	querybuilder "homeserver/db/query_builder"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
	"slices"
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

	query := querybuilder.New().
		Update(dbconsts.DBFilesTable.Name).
		Set(fields_to_update...).
		Where("id = $" + fmt.Sprint(len(fields_to_update)+1)).
		Returning(fschema.FieldsList()...).
		Build()

	// query := fmt.Sprintf(`
	// 	UPDATE user_files
	// 	SET %s
	// 	WHERE id = $%d
	// 	RETURNING`+strings.Join(fschema.FieldsList(), ", "),
	// 	fields_str, len(fields_to_update)+1)

	query_args := make([]any, len(fields_to_update)+1)
	copy(query_args, OrderFileUpdateFields(fields_to_update, modified_file))
	query_args[len(fields_to_update)] = modified_file.ID

	rows, err := repo.db.Pool.Query(ctx, query, query_args...)
	if err != nil {
		return fmt.Errorf("error updating file entry: %v", err)
	}

	var modified_files []fschema.UserFile
	if err := helpers.ScanRows(rows, &modified_files); err != nil {
		return fmt.Errorf("error scanning updated file entry: %v", err)
	}

	if len(modified_files) != 1 {
		return fmt.Errorf("expected 1 modified file, got %d", len(modified_files))
	}

	*modified_file = modified_files[0]

	return nil
}
