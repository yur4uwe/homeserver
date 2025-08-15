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

func (repo *FileRepo) GetFileByID(ctx context.Context, id string) (*fschema.UserFile, error) {
	fields_to_extract := strings.Join([]string{
		fconsts.ID, fconsts.Filename, fconsts.Extension,
		fconsts.Path, fconsts.Directory, fconsts.Size,
		fconsts.ContentType, fconsts.OwnerID, fconsts.CreatedAt,
		fconsts.UpdatedAt, fconsts.Metadata,
	}, ", ")

	var file fschema.UserFile
	err := repo.db.Pool.QueryRow(
		ctx,
		"SELECT "+fields_to_extract+" FROM files WHERE id = $1",
		id,
	).Scan(
		&file.ID, &file.Filename, &file.Extension,
		&file.Path, &file.Directory, &file.Size,
		&file.ContentType, &file.OwnerID, &file.CreatedAt,
		&file.UpdatedAt, &file.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (repo *FileRepo) GetFullFileInfo(ctx context.Context, file *fschema.UserFile, owner *uschema.User) (*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.ID, uconsts.Username); err != nil {
		return nil, err
	}

	if err := file.ShouldHaveFields(owner.Username)(fschema.FieldsList()...); err != nil {
		return nil, err
	}

	var fullFileInfo fschema.UserFile
	err := repo.db.Pool.QueryRow(
		ctx,
		"SELECT "+strings.Join(fschema.FieldsList(), ", ")+" FROM files WHERE filename = $1 AND path = $2 AND owner_id = $3",
		file.Filename, file.Path, owner.ID,
	).Scan(
		&fullFileInfo.ID, &fullFileInfo.Filename, &fullFileInfo.Extension,
		&fullFileInfo.Path, &fullFileInfo.Directory, &fullFileInfo.Size,
		&fullFileInfo.ContentType, &fullFileInfo.OwnerID, &fullFileInfo.CreatedAt,
		&fullFileInfo.UpdatedAt, &fullFileInfo.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return &fullFileInfo, nil
}

func (repo *FileRepo) GetFilesByUser(ctx context.Context, owner *uschema.User) ([]*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.ID); err != nil {
		return nil, err
	}

	var files []*fschema.UserFile
	rows, err := repo.db.Pool.Query(
		ctx,
		"SELECT "+strings.Join(fschema.FieldsList(), ", ")+" FROM files WHERE owner_id = $1",
		owner.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (repo *FileRepo) FuzzySearch(ctx context.Context, search_query string, directory string, owner *uschema.User) ([]*fschema.UserFile, error) {
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

	var files []*fschema.UserFile
	rows, err := repo.db.Pool.Query(ctx, query, search_query, owner.ID)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}
	defer rows.Close()

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

		if strings.HasPrefix(file.Path, directory) {
			files = append(files, &file)
		}
	}

	return files, nil
}

func (repo *FileRepo) GetFilesByDirectory(ctx context.Context, directory string, owner *uschema.User) ([]*fschema.UserFile, error) {
	if err := owner.ShouldHaveFields(uconsts.ID); err != nil {
		return nil, fmt.Errorf("owner must have a valid user ID")
	}

	fields := fschema.FieldsList()
	args := make([]any, len(fields))
	for i, f := range fields {
		args[i] = f
	}

	query := fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s
		FROM user_files
		WHERE path = $1 AND owner_id = $2
	`, args...)

	rows, err := repo.db.Pool.Query(ctx, query, directory, owner.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying files in directory: %v", err)
	}
	defer rows.Close()

	var files []*fschema.UserFile
	for rows.Next() {
		var file fschema.UserFile
		if err := rows.Scan(&file.ID, &file.Filename, &file.Extension, &file.Path, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt, &file.UpdatedAt, &file.Metadata); err != nil {
			return nil, fmt.Errorf("error scanning file entry: %v", err)
		}

		if strings.HasPrefix(file.Path, directory) {
			files = append(files, &file)
		}
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
	defer rows.Close()

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
