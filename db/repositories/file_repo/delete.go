package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/fconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbdelete"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
)

func (repo *FileRepo) DeleteFileByID(ctx context.Context, fileID int, owner *uschema.User) (*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.ID); err != nil {
		return nil, err
	}

	query := qbdelete.New().From(dbconsts.DBFilesTable.Name).Where("id = $1 AND owner_id = $2").Returning(fschema.FieldsList()...).Build()

	deleted_file, err := repo.db.QueryOne(ctx, query, fileID, owner.ID)
	if err != nil {
		return nil, err
	}

	return &deleted_file, nil
}

func (repo *FileRepo) DeleteFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.Username, uconsts.ID); err != nil {
		return nil, err
	}

	if err := file.ShouldHaveFields(owner.Username)(fconsts.Filename, fconsts.Path); err != nil {
		return nil, err
	}

	whereClause := ""
	for i, field := range file.GetValidFieldsList() {
		if i > 0 {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%s = $%d", field, i+1)
	}

	args := make([]any, 0, len(file.GetValidFieldsList()))
	args = append(args, file.GetValidFieldsData()...)

	query := qbdelete.New().From(dbconsts.DBFilesTable.Name).Where(whereClause).Returning(fschema.FieldsList()...).Build()

	deleted_file, err := repo.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &deleted_file, nil
}
