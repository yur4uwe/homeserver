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

func (repo *FileRepo) UpdateFileEntry(ctx context.Context, file *fschema.UserFile, fvp map[string]any, owner *uschema.User) error {
	keys := []string{}
	fval := *file
	projected_file := &fval

	for k := range fvp {
		if _, err := file.Get(k); err == nil {
			keys = append(keys, k)
			projected_file.Set(k, fvp[k])
		}
	}

	if slices.Contains(keys, fconsts.ID) {
		return fmt.Errorf("cannot update file ID")
	}
	if slices.Contains(keys, fconsts.OwnerID) {
		return fmt.Errorf("cannot update owner ID")
	}
	if slices.Contains(keys, fconsts.CreatedAt) {
		return fmt.Errorf("cannot update created at")
	}
	if slices.Contains(keys, fconsts.UpdatedAt) {
		return fmt.Errorf("cannot update updated at")
	}

	if err := file.ShouldHaveFields(owner.Username)(fconsts.ID); err != nil {
		return fmt.Errorf("invalid file data: %v", err)
	}

	var err error
	if slices.Contains(keys, fconsts.Filename) || slices.Contains(keys, fconsts.Extension) || slices.Contains(keys, fconsts.Path) {
		projected_file, err = FixDuplicates(projected_file)
		if err != nil {
			return fmt.Errorf("error preventing duplicate file name: %v", err)
		}

		fvp[fconsts.Filename] = projected_file.Filename
	}

	query := querybuilder.New().
		Update(dbconsts.DBFilesTable.Name).
		Set(keys...).
		Where("id = $" + fmt.Sprint(len(keys)+1)).
		Returning(fschema.FieldsList()...).
		Build()

	// query := fmt.Sprintf(`
	// 	UPDATE user_files
	// 	SET %s
	// 	WHERE id = $%d
	// 	RETURNING`+strings.Join(fschema.FieldsList(), ", "),
	// 	fields_str, len(fields_to_update)+1)

	query_args := make([]any, 0, len(keys)+1)
	for _, k := range keys {
		query_args = append(query_args, fvp[k])
	}
	query_args = append(query_args, file.ID)

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

	*file = modified_files[0]

	return nil
}
