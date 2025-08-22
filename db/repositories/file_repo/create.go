package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbinsert"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (repo *FileRepo) CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.Username); err != nil {
		return nil, err
	}

	fields := fschema.FieldsList(fconsts.ID, fconsts.CreatedAt, fconsts.UpdatedAt, fconsts.Metadata)

	if err := file.ShouldHaveFields(owner.Username)(fields...); err != nil {
		return nil, err
	}

	file, err := PreventDuplicateFileName(file)
	if err != nil {
		return nil, err
	}

	query := qbinsert.New().
		Columns(fields...).
		Into(dbconsts.DBFilesTable.Name).
		Returning(fschema.FieldsList()...).
		Build()

	var inserted_file []fschema.UserFile
	rows, err := repo.db.Pool.Query(
		ctx,
		query,
		file.ID, file.Filename, file.Extension, file.Path, file.Directory, file.Size, file.ContentType, file.OwnerID, file.Metadata,
	)
	if err != nil {
		return nil, err
	}

	if err = helpers.ScanRows(rows, &inserted_file); err != nil {
		return nil, err
	}

	if len(inserted_file) != 1 {
		return nil, fmt.Errorf("unexpected error: %d rows returned", len(inserted_file))
	}

	return &inserted_file[0], nil
}
