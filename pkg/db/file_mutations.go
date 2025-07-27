package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"strings"

	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func PreventDuplicateFileName(filePath string, owner string) (string, error) {
	// Extract the relative path and filename
	storagePath := os.Getenv(globals.EnvStoragePath)

	relativePath := strings.ReplaceAll(filePath, "\\", "/")
	relativePath = strings.TrimPrefix(relativePath, storagePath+owner)
	relativePath = strings.TrimPrefix(relativePath, "/") // Remove leading slash if present

	fileName := filepath.Base(relativePath)
	if fileName == "." || fileName == "/" {
		return "", fmt.Errorf("invalid file path: %s", filePath)
	}

	dirPath := filepath.Dir(relativePath)
	if dirPath == "." {
		dirPath = ""
	}

	base := fileName
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	counter := 1

	newName := fmt.Sprintf("%s%s", name, ext)
	newPath := BuildFilePath(owner, dirPath, newName)

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, nil
		}

		newName = fmt.Sprintf("%s_%d%s", name, counter, ext)
		newPath = BuildFilePath(owner, dirPath, newName)

		counter++
	}
}

func BuildFilePath(owner string, relative_path string, filename string) string {
	storage_path := os.Getenv(globals.EnvStoragePath)
	if storage_path[len(storage_path)-1] != '/' {
		panic("Storage path must end with a slash")
	}

	var file_path string
	if relative_path != "" {
		file_path = fmt.Sprintf("/%s/%s", relative_path, filename)
	} else {
		file_path = fmt.Sprintf("/%s", filename)
	}

	return fmt.Sprintf("%s%s%s", storage_path, owner, file_path)
}

func (db *Database) SaveFile(
	ctx context.Context,
	header *UserFile,
	file *multipart.File,
	requester *User,
) (int, string, error) {
	// Define the file path
	filePath := BuildFilePath(requester.Username, header.Path, header.Filename)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return 0, "", fmt.Errorf("error creating directory: %v", err)
	}

	// Check if the file already exists
	filePath, err := PreventDuplicateFileName(filePath, requester.Username)
	if err != nil {
		return 0, "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	header.Filename = filepath.Base(filePath)

	// Save the file to disk
	outFile, err := os.Create(filePath)
	if err != nil {
		return 0, "", fmt.Errorf("error saving file to disk: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, *file)
	if err != nil {
		return 0, "", fmt.Errorf("error writing file to disk: %v", err)
	}

	filePath = strings.ReplaceAll(filePath, "\\", "/")

	// Insert file metadata into the database
	var fileID int
	query := fmt.Sprintf(`
		INSERT INTO user_files (%s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, Filename, Path, OwnerID, Size, ContentType)

	err = db.pool.QueryRow(ctx, query, header.Filename, filePath, requester.ID, header.Size, header.ContentType).Scan(&fileID)
	if err != nil {
		return 0, "", fmt.Errorf("error inserting file metadata: %v", err)
	}

	return fileID, header.Filename, nil
}

func (db *Database) DeleteFile(ctx context.Context, filePath string, owner_token string) error {
	user, err := db.QueryToken(ctx, owner_token)
	if err != nil {
		return fmt.Errorf("error querying token: %v", err)
	}

	filename := filepath.Base(filePath)
	if filename == "." || filename == "/" {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// Get the file path from the database
	path := BuildFilePath(user.Username, "", filePath)

	// Delete the file from disk
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error deleting file from disk: %v", err)
	}

	// Delete the file metadata from the database
	_, err = db.pool.Exec(ctx, `DELETE FROM user_files WHERE filename = $1 AND owner_id = $2`, filename, user.ID)
	if err != nil {
		return fmt.Errorf("error deleting file metadata: %v", err)
	}

	return nil
}

func (db *Database) RenameFile(ctx context.Context, oldPath string, newPath string, owner_token string) (string, error) {
	user, err := db.QueryToken(ctx, owner_token)
	if err != nil {
		return "", fmt.Errorf("error querying token: %v", err)
	}

	oldFilePath := BuildFilePath(user.Username, "", oldPath)
	newFilePath := BuildFilePath(user.Username, "", newPath)

	newFilePath, err = PreventDuplicateFileName(newFilePath, user.Username)
	if err != nil {
		return "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	newFileName := filepath.Base(newFilePath)

	if oldFilePath == newFilePath {
		return newFilePath, nil // No change needed
	}

	// Rename the file on disk
	if err := os.Rename(oldFilePath, newFilePath); err != nil {
		return "", fmt.Errorf("error renaming file on disk: %v", err)
	}

	// Update the file metadata in the database
	_, err = db.pool.Exec(ctx, `UPDATE user_files SET path = $1, filename = $2 WHERE path = $3 AND owner_id = $4`, newFilePath, newFileName, oldFilePath, user.ID)
	if err != nil {
		return "", fmt.Errorf("error updating file metadata: %v", err)
	}

	return newFilePath, nil
}
