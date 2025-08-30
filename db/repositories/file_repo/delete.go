package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbdelete"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (repo *FileRepo) DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.Username, uconsts.ID); err != nil {
		return nil, err
	}

	identifier, args := file.Identificator(owner.Username, 1)
	if identifier == "" {
		return nil, fmt.Errorf("file is missing required fields")
	}

	query := qbdelete.New().
		From(dbconsts.DBFilesTable.Name).
		Where(identifier).
		Returning(fschema.FieldsList()...).
		Build()

	deleted_file, err := repo.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &deleted_file, nil
}
