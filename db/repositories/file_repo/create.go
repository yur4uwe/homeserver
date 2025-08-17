package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/fconsts"
	"homeserver/constants/uconsts"
	fschema "homeserver/db/repositories/file_repo/schema"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strings"
)

func (repo *FileRepo) CreateFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.Username); err != nil {
		return nil, err
	}

	fields := []string{
		fconsts.Filename, fconsts.Extension,
		fconsts.Path, fconsts.Directory, fconsts.Size,
		fconsts.ContentType, fconsts.OwnerID,
	}

	if err := file.ShouldHaveFields(owner.Username)(fields...); err != nil {
		return nil, err
	}

	file, err := PreventDuplicateFileName(file)
	if err != nil {
		return nil, err
	}

	fields_to_extract := strings.Join(fields, ", ")
	insert_points := ""
	for i := 1; i <= len(fields); i++ {
		insert_points += fmt.Sprintf("$%d", i)
		if i < len(fields) {
			insert_points += ", "
		}
	}

	var inserted_file fschema.UserFile
	err = repo.db.Pool.QueryRow(
		ctx,
		"INSERT INTO files ("+fields_to_extract+") VALUES ("+insert_points+") RETURNING "+fields_to_extract,
		file.ID, file.Filename, file.Extension, file.Path, file.Directory, file.Size, file.ContentType, file.OwnerID, file.Metadata,
	).Scan(
		&inserted_file.ID, &inserted_file.Filename, &inserted_file.Extension,
		&inserted_file.Path, &inserted_file.Directory, &inserted_file.Size,
		&inserted_file.ContentType, &inserted_file.OwnerID, &inserted_file.CreatedAt,
		&inserted_file.UpdatedAt, &inserted_file.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return &inserted_file, nil
}
