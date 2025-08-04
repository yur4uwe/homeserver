package file

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"strings"
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

func SearchFiles(ctx context.Context, db interfaces.DatabaseInterface, query string, requester *user_schema.User, directory []string) ([]*file_schema.UserFile, error) {
	files, err := ExtractFilesFromDirectory(ctx, db, strings.Join(directory, "/"), requester)
	if err != nil {
		return nil, fmt.Errorf("error extracting user files: %v", err)
	}

	return files, nil
}

func ExtractFilesFromDirectory(ctx context.Context, db interfaces.DatabaseInterface, directory string, requester *user_schema.User) ([]*file_schema.UserFile, error) {
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

func GetFile(ctx context.Context, db interfaces.DatabaseInterface, filePath string, requester *user_schema.User) (*file_schema.UserFile, error) {
	query := fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s, %s, %s
		FROM user_files
		WHERE path = $1 AND owner_id = $2
	`, schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Size, schemas.ContentType, schemas.OwnerID, schemas.CreatedAt)

	row := db.QueryRow(ctx, query, filePath, requester.ID)

	file := &file_schema.UserFile{}
	if err := row.Scan(&file.ID, &file.Filename, &file.Extension, &file.Path, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt); err != nil {
		return nil, fmt.Errorf("error retrieving file: %v", err)
	}

	return file, nil
}
