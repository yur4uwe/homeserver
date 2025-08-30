package filerepo

import (
	"context"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbinsert"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (repo *FileRepo) CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.Username); err != nil {
		return nil, err
	}

	fields := fschema.FieldsList(fconsts.ID, fconsts.CreatedAt, fconsts.UpdatedAt, fconsts.Metadata, fconsts.TrustedSource)

	if err := file.ShouldHaveFields(owner.Username)(fields...); err != nil {
		return nil, err
	}

	file, err := FixDuplicates(file)
	if err != nil {
		return nil, err
	}

	args := file.GetValidFieldsData()

	query := qbinsert.New().
		Columns(fields...).
		Into(dbconsts.DBFilesTable.Name).
		Returning(fschema.FieldsList()...).
		Build()

	inserted_file, err := repo.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &inserted_file, nil
}
