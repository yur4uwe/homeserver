package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"homeserver/db/query_builder/qbupdate"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
)

var (
	FixDuplicates = fschema.PreventDuplicateFileName
	GetFileByID   = (*FileRepo).GetFileByID
)

func (repo *FileRepo) UpdateFileEntry(ctx context.Context, file *fschema.UserFile, fvp *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := fvp.ShouldNotHaveAnyOf(owner.Username)(fconsts.ID, fconsts.OwnerID, fconsts.CreatedAt, fconsts.UpdatedAt); err != nil {
		return nil, fmt.Errorf("inappropriate fields found: %v", err)
	}

	if err := file.ShouldHaveFields(owner.Username)(fconsts.ID); err != nil {
		return nil, fmt.Errorf("invalid file data: %v", err)
	}

	file, err := GetFileByID(repo, ctx, file.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting full file info: %v", err)
	}

	file.MergeFiles(fvp)

	query := qbupdate.New().
		Set(file.GetValidFieldsList()...).
		From(dbconsts.DBFilesTable.Name).
		Where("id = $" + fmt.Sprint(len(file.GetValidFieldsList())+1)).
		Returning(fschema.FieldsList()...).
		Build()

	query_args := make([]any, 0, len(file.GetValidFieldsList())+1)
	query_args = append(query_args, file.GetValidFieldsData()...)
	query_args = append(query_args, file.ID)

	rows, err := repo.db.Pool.Query(ctx, query, query_args...)
	if err != nil {
		return nil, fmt.Errorf("error updating file entry: %v", err)
	}

	var modified_files []fschema.UserFile
	if err := helpers.ScanRows(rows, &modified_files); err != nil {
		return nil, fmt.Errorf("error scanning updated file entry: %v", err)
	}

	if len(modified_files) != 1 {
		return nil, fmt.Errorf("expected 1 modified file, got %d", len(modified_files))
	}

	*file = modified_files[0]

	return file, nil
}
