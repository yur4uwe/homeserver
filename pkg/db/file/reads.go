package file

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
)

func ExtractUserFiles(ctx context.Context, db interfaces.DatabaseInterface, userID int) ([]*file_schema.UserFile, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	rows, err := db.Query(
		ctx,
		fmt.Sprintf(`SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s, %s FROM user_files WHERE %s = $1`,
			schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.CreatedAt, schemas.UpdatedAt, schemas.Metadata, schemas.OwnerID,
		),
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying user files: %v", err)
	}
	defer rows.Close()

	var entries []*file_schema.UserFile
	for rows.Next() {
		entry := &file_schema.UserFile{}
		if err := rows.Scan(&entry.ID, &entry.Filename, &entry.Extension, &entry.Path, &entry.Size, &entry.ContentType, &entry.OwnerID, &entry.CreatedAt, &entry.UpdatedAt, &entry.Metadata); err != nil {
			return nil, fmt.Errorf("error scanning user file entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func SearchFiles(ctx context.Context, db interfaces.DatabaseInterface, search_query string, requester *user_schema.User, directory []string) ([]*file_schema.UserFile, error) {
	select_query := fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s, %s, %s FROM user_files`,
		schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path,
		schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.CreatedAt,
	)

	query := select_query + `
		WHERE filename % $1
		ORDER BY similarity(filename, $1) DESC
		LIMIT 20;`

	rows, err := db.Query(ctx, query, search_query)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}
	defer rows.Close()

	var files []*file_schema.UserFile
	for rows.Next() {
		file := &file_schema.UserFile{}
		if err := rows.Scan(
			&file.ID, &file.Filename, &file.Extension, &file.Path,
			&file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning file entry: %v", err)
		}
		files = append(files, file)
	}

	if !rows.Next() && rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", rows.Err())
	}

	return files, nil
}

func ExtractFilesFromDirectory(ctx context.Context, db interfaces.DatabaseInterface, directory string, requester *user_schema.User) ([]*file_schema.UserFile, error) {
	if err := requester.ShouldHaveFields(schemas.UserID); err != nil {
		return nil, fmt.Errorf("requester must have a valid user ID")
	}

	query := fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s, %s, %s
		FROM user_files
		WHERE path = $1 AND owner_id = $2
	`, schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.CreatedAt)

	rows, err := db.Query(ctx, query, directory, requester.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying files in directory: %v", err)
	}
	defer rows.Close()

	var files []*file_schema.UserFile
	for rows.Next() {
		file := &file_schema.UserFile{}
		if err := rows.Scan(&file.ID, &file.Filename, &file.Extension, &file.Path, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning file entry: %v", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func ExtractAllFiles(ctx context.Context, db interfaces.DatabaseInterface) ([]*file_schema.UserFile, error) {
	rows, err := db.Query(
		ctx,
		fmt.Sprintf(`SELECT %s, %s, %s, %s, %s, %s, %s FROM user_files`,
			schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Size, schemas.ContentType, schemas.OwnerID,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error querying all files: %v", err)
	}
	defer rows.Close()

	var files []*file_schema.UserFile
	for rows.Next() {
		file := &file_schema.UserFile{}
		if err := rows.Scan(&file.ID, &file.Filename, &file.Extension, &file.Path, &file.Size, &file.ContentType, &file.OwnerID); err != nil {
			return nil, fmt.Errorf("error scanning file entry: %v", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func GetFile(ctx context.Context, db interfaces.DatabaseInterface, file *file_schema.UserFile, requester *user_schema.User) (*file_schema.UserFile, error) {
	if err := requester.ShouldHaveFields(schemas.Username); err != nil {
		return nil, fmt.Errorf("error validating requester: %v", err)
	}

	if err := file.ShouldHaveFields(requester.Username)(schemas.Filename, schemas.Path); err != nil {
		return nil, fmt.Errorf("invalid file data: %v", err)
	}

	rows, err := db.Query(ctx, fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s, %s, %s
		FROM user_files
	`, schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.CreatedAt))

	if err != nil {
		return nil, fmt.Errorf("error querying file: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&file.ID, &file.Filename, &file.Extension, &file.Path, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning file entry: %v", err)
		}
	}

	query := fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s, %s, %s
		FROM user_files
		WHERE path = $1 AND filename = $2 AND owner_id = $3
	`, schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.CreatedAt)

	row := db.QueryRow(ctx, query, file.Path, file.Filename, requester.ID)

	res_file := &file_schema.UserFile{}
	if err := row.Scan(&res_file.ID, &res_file.Filename, &res_file.Extension, &res_file.Path, &res_file.Size, &res_file.ContentType, &res_file.OwnerID, &res_file.CreatedAt); err != nil {
		return nil, fmt.Errorf("error retrieving res_file: %v", err)
	}

	return res_file, nil
}
