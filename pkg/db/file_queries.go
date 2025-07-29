package db

import (
	"context"
	"fmt"
	"strings"
)

func (db *Database) ExtractUserFiles(ctx context.Context, userID int) ([]*UserFile, error) {
	rows, err := db.pool.Query(
		ctx,
		fmt.Sprintf(`SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s, %s FROM user_files WHERE %s = $1`,
			FileID, Filename, Extension, Path, Size, ContentType, OwnerID, CreatedAt, UpdatedAt, Metadata, OwnerID,
		),
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying user files: %v", err)
	}
	defer rows.Close()

	var entries []*UserFile
	for rows.Next() {
		entry := &UserFile{}
		if err := rows.Scan(&entry.ID, &entry.Filename, &entry.Extension, &entry.Path, &entry.Size, &entry.ContentType, &entry.OwnerID, &entry.CreatedAt, &entry.UpdatedAt, &entry.Metadata); err != nil {
			return nil, fmt.Errorf("error scanning user file entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (db *Database) SearchFiles(ctx context.Context, query string, requester *User, directory []string) ([]*UserFile, error) {
	files, err := db.ExtractFilesFromDirectory(ctx, strings.Join(directory, "/"), requester)
	if err != nil {
		return nil, fmt.Errorf("error extracting user files: %v", err)
	}

	return files, nil
}

func (db *Database) CreateFileEntry(ctx context.Context, file *UserFile, requester *User) (int, error) {
	if file.Filename == "" || file.Size == 0 || file.ContentType == "" || requester == nil {
		return -1, fmt.Errorf("invalid file entry data")
	}

	if strings.Contains(file.Path, "..") || strings.Contains(file.Path, "//") || strings.Contains(file.Path, "\\") {
		return -1, fmt.Errorf("invalid file path: %s", file.Path)
	}

	relative_path := strings.Split(file.Path, "/")
	if len(relative_path) == 0 || relative_path[0] == "" {
		relative_path = []string{requester.Username}
	}

	id, err := db.QueryFolderByName(ctx, relative_path[len(relative_path)-1], requester.ID, strings.Join(relative_path[:len(relative_path)-1], "/"))
	if err != nil {
		return -1, fmt.Errorf("error querying folder by name: %v", err)
	}
	file.Directory = id

	var fileID int
	query := fmt.Sprintf(`
		INSERT INTO user_files (%s, %s, %s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, Filename, Extension, Path, OwnerID, Size, ContentType, Directory)

	err = db.pool.QueryRow(ctx, query, file.Filename, file.Extension, strings.Join(relative_path, "/"), requester.ID, file.Size, file.ContentType, file.Directory).Scan(&fileID)
	if err != nil {
		return -1, fmt.Errorf("error inserting file metadata: %v", err)
	}

	return fileID, nil
}

func (db *Database) ExtractFilesFromDirectory(ctx context.Context, directory string, requester *User) ([]*UserFile, error) {
	query := fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s, %s, %s
		FROM user_files
		WHERE path = $1 AND owner_id = $2
	`, FileID, Filename, Extension, Path, Size, ContentType, OwnerID, CreatedAt)

	rows, err := db.pool.Query(ctx, query, directory, requester.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying files in directory: %v", err)
	}
	defer rows.Close()

	var files []*UserFile
	for rows.Next() {
		file := &UserFile{}
		if err := rows.Scan(&file.ID, &file.Filename, &file.Extension, &file.Path, &file.Size, &file.ContentType, &file.OwnerID, &file.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning file entry: %v", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func (db *Database) ExtractAllFiles(ctx context.Context) ([]*UserFile, error) {
	rows, err := db.pool.Query(ctx, fmt.Sprintf(`SELECT %s, %s, %s, %s, %s, %s, %s FROM user_files`, FileID, Filename, Extension, Path, Size, ContentType, OwnerID))
	if err != nil {
		return nil, fmt.Errorf("error querying all files: %v", err)
	}
	defer rows.Close()

	var files []*UserFile
	for rows.Next() {
		file := &UserFile{}
		if err := rows.Scan(&file.ID, &file.Filename, &file.Extension, &file.Path, &file.Size, &file.ContentType, &file.OwnerID); err != nil {
			return nil, fmt.Errorf("error scanning file entry: %v", err)
		}
		files = append(files, file)
	}

	return files, nil
}
