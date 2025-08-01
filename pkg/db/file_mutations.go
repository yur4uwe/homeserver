package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"strings"

	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func (db *Database) SaveFile(
	ctx context.Context,
	header *file_schema.UserFile,
	file *multipart.File,
	requester *user_schema.User,
) (int, string, error) {
	filePath := globals.BuildFilePath(header.Path, header.Filename, header.Extension)

	path_slice := strings.Split(header.Path, "/")
	if err := EnsureFolderExists(ctx, db, path_slice, requester.ID); err != nil {
		return 0, "", fmt.Errorf("error ensuring parent folders exist: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return 0, "", fmt.Errorf("error creating directory: %v", err)
	}

	filePath, err := PreventDuplicateFileName(filePath)
	if err != nil {
		return 0, "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	header.Filename = strings.Split(filepath.Base(filePath), ".")[0]

	outFile, err := os.Create(filePath)
	if err != nil {
		return 0, "", fmt.Errorf("error saving file to disk: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, *file)
	if err != nil {
		return 0, "", fmt.Errorf("error writing file to disk: %v", err)
	}

	fileID, err := db.CreateFileEntry(ctx, header, requester)
	if err != nil {
		return 0, "", fmt.Errorf("error creating file entry: %v", err)
	}

	return fileID, header.Filename + "." + header.Extension, nil
}

func (db *Database) DeleteFile(ctx context.Context, file *file_schema.UserFile, user_name string) (int, error) {

	file_path_slice := strings.Split(file.Path, "/")
	if len(file_path_slice) < 2 {
		return 0, fmt.Errorf("invalid file path: %s", file.Path)
	}

	file_path_slice = file_path_slice[:len(file_path_slice)-1]

	path := globals.BuildFilePath(strings.Join(file_path_slice, "/"), file.Filename, file.Extension)

	if err := os.Remove(path); err != nil {
		return 0, fmt.Errorf("error deleting file from disk: %v", err)
	}

	var id int
	err := db.pool.QueryRow(
		ctx,
		`DELETE FROM user_files WHERE filename = $1 AND owner_id = $2 RETURNING id`,
		file.Filename, file.OwnerID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error deleting file metadata: %v", err)
	}

	return id, nil
}

func (db *Database) RenameFile(ctx context.Context, oldPath string, newPath string, requester *user_schema.User) (string, error) {
	oldRelFilePath := requester.Username + "/" + oldPath
	newRelFilePath := requester.Username + "/" + newPath

	oldAbsFilePath := globals.BuildFilePath(requester.Username, strings.Split(oldPath, ".")[0], strings.Split(oldPath, ".")[1])

	newAbsFilePath, err := PreventDuplicateFileName(newRelFilePath)
	if err != nil {
		return "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	newFileName := strings.Split(filepath.Base(newAbsFilePath), ".")[0]
	oldFileName := strings.Split(filepath.Base(oldAbsFilePath), ".")[0]

	if oldRelFilePath == newRelFilePath {
		return newRelFilePath, nil
	}

	if err := os.Rename(oldAbsFilePath, newAbsFilePath); err != nil {
		return "", fmt.Errorf("error renaming file on disk: %v", err)
	}

	relPath := strings.Join(strings.Split(newRelFilePath, "/")[:len(strings.Split(newRelFilePath, "/"))-1], "/")

	tag, err := db.pool.Exec(
		ctx,
		`UPDATE user_files SET path = $1, filename = $2, updated_at = NOW() WHERE filename = $3 AND path = $4 AND owner_id = $5`,
		relPath, newFileName, oldFileName, relPath, requester.ID,
	)
	if err != nil {
		return "", fmt.Errorf("error updating file metadata: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return "", fmt.Errorf("no file found to rename: %s", oldRelFilePath)
	}

	return newRelFilePath, nil
}

func (db *Database) UpdateFileEntry(ctx context.Context, file *file_schema.UserFile) error {
	if file.Filename == "" || file.Size == 0 || file.ContentType == "" {
		return fmt.Errorf("invalid file entry data")
	}

	query := `
		UPDATE user_files
		SET filename = $1, extension = $2, path = $3, size = $4, content_type = $5, updated_at = NOW()
		WHERE id = $6
	`

	_, err := db.pool.Exec(ctx, query, file.Filename, file.Extension, file.Path, file.Size, file.ContentType, file.ID)
	if err != nil {
		return fmt.Errorf("error updating file entry: %v", err)
	}

	return nil
}
