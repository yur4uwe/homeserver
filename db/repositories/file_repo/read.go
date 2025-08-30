package filerepo

import (
	"context"
	"fmt"
	"homeserver/constants/dbconsts"
	"homeserver/constants/uconsts"
	"homeserver/db/query_builder/qbselect"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/db/repositories/helpers"
	uschema "homeserver/db/repositories/user_repo/schema"
	"strings"
)

func (repo *FileRepo) GetFile(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	uname := ""
	if owner != nil {
		uname = owner.Username
	}

	identifier, args := file.Identificator(uname, 1)
	if identifier == "" {
		return nil, fmt.Errorf("file is missing required fields")
	}

	query := qbselect.New().
		Columns(fschema.FieldsList()...).
		From(dbconsts.DBFilesTable.Name).
		Where(identifier).
		Build()

	extracted_file, err := repo.db.QueryOne(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &extracted_file, nil
}

func (repo *FileRepo) GetFilesByUser(ctx context.Context, owner *uschema.User) ([]fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.ID); err != nil {
		return nil, err
	}

	query := qbselect.New().
		Columns(fschema.FieldsList()...).
		From(dbconsts.DBFilesTable.Name).
		Where("owner_id = $1").
		Build()
	// query := querybuilder.New().
	// 	Select(fschema.FieldsList()...).
	// 	From(dbconsts.DBFilesTable.Name).
	// 	Where("owner_id = $1").
	// 	Build()

	var files []fschema.UserFile
	rows, err := repo.db.Pool.Query(
		ctx,
		query,
		owner.ID,
	)
	if err != nil {
		return nil, err
	}

	if err := helpers.ScanRows(rows, &files); err != nil {
		return nil, fmt.Errorf("error scanning files: %v", err)
	}

	return files, nil
}

func (repo *FileRepo) FuzzySearch(ctx context.Context, search_query string, directory string, owner *uschema.User) ([]fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.ID); err != nil {
		return nil, err
	}

	fields := strings.Join(fschema.FieldsList(), ", f.")
	select_query := fmt.Sprintf(
		`SELECT f.%s
		FROM user_files f
		INNER JOIN users u ON f.owner_id = u.id AND u.id = $2`,
		fields,
	)

	query := select_query + `
	    WHERE f.filename % $1
	    ORDER BY similarity(f.filename, $1) DESC
	    LIMIT 20;`

	var files []fschema.UserFile
	rows, err := repo.db.Pool.Query(ctx, query, search_query, owner.ID)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}

	if err := helpers.ScanRows(rows, &files); err != nil {
		return nil, fmt.Errorf("error scanning files: %v", err)
	}

	return files, nil
}

func (repo *FileRepo) GetFilesByDirectory(ctx context.Context, directory string, owner *uschema.User) ([]fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.ID); err != nil {
		return nil, fmt.Errorf("owner must have a valid user ID")
	}

	fields := fschema.FieldsList()
	args := make([]any, len(fields))
	for i, f := range fields {
		args[i] = f
	}

	builder := qbselect.New().
		From(dbconsts.DBFilesTable.Name).
		Where("path = $1 AND owner_id = $2")

	for _, field := range fields {
		builder.Col(field)
	}

	query := builder.Build()

	rows, err := repo.db.Pool.Query(ctx, query, directory, owner.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying files in directory: %v", err)
	}

	var files []fschema.UserFile
	if err := helpers.ScanRows(rows, &files); err != nil {
		return nil, fmt.Errorf("error scanning files: %v", err)
	}

	return files, nil
}

func (repo *FileRepo) GetFiles(ctx context.Context) ([]*fschema.UserFile, error) {
	var files []*fschema.UserFile
	rows, err := repo.db.Pool.Query(
		ctx,
		"SELECT "+strings.Join(fschema.FieldsList(), ", ")+" FROM files",
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var file fschema.UserFile
		if err := rows.Scan(
			&file.ID, &file.Filename, &file.Extension,
			&file.Path, &file.Directory, &file.Size,
			&file.ContentType, &file.OwnerID, &file.CreatedAt,
			&file.UpdatedAt, &file.Metadata,
		); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	return files, nil
}
