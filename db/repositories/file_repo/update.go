package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"homeserver/db/query_builder/qbupdate"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

var (
	FixDuplicates = fschema.PreventDuplicateFileName
)

func (repo *FileRepo) UpdateFileEntry(ctx context.Context, file *fschema.UserFile, fvp *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := fvp.ShouldNotHaveAnyOf(owner.Username)(fconsts.ID, fconsts.OwnerID, fconsts.CreatedAt, fconsts.UpdatedAt); err != nil {
		return nil, fmt.Errorf("inappropriate fields found: %v", err)
	}

	if err := file.ShouldHaveFields(owner.Username)(fconsts.ID); err != nil {
		return nil, fmt.Errorf("invalid file data: %v", err)
	}

	fields_to_update := fvp.GetValidFieldsList()

	query := qbupdate.New().
		Set(fields_to_update...).
		From(dbconsts.DBFilesTable.Name).
		Where(fmt.Sprintf("id = $%d", len(fields_to_update)+1)).
		Returning(fschema.FieldsList()...).
		Build()

	query_args := make([]any, 0, len(fields_to_update)+1)
	query_args = append(query_args, fvp.GetValidFieldsData()...)
	query_args = append(query_args, file.ID)

	modded_file, err := repo.db.QueryOne(ctx, query, query_args...)
	if err != nil {
		return nil, fmt.Errorf("error updating file entry: %v", err)
	}

	return &modded_file, nil
}
