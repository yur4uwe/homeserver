package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/templates"

	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func BuildFilePath(owner string, path string, filename string) string {
	storage_path := os.Getenv(globals.EnvStoragePath)
	if storage_path[len(storage_path)-1] != '/' {
		panic("Storage path must end with a slash")
	}

	var file_path string
	if path != "" {
		file_path = fmt.Sprintf("/%s/%s", path, filename)
	} else {
		file_path = fmt.Sprintf("/%s", filename)
	}

	return fmt.Sprintf("%s%s%s", storage_path, owner, file_path)
}

func (db *Database) SaveFile(
	ctx context.Context,
	header *templates.FileInfo,
	file *multipart.File,
	token string,
) (int, error) {
	user, err := db.QueryToken(ctx, token)
	if err != nil {
		return 0, fmt.Errorf("error querying token: %v", err)
	}

	// Define the file path
	filePath := BuildFilePath(user.Username, header.Path, header.Name)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return 0, fmt.Errorf("error creating directory: %v", err)
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		// If the file exists, append a number to the filename
		base := header.Name
		ext := filepath.Ext(base)
		name := base[:len(base)-len(ext)]
		counter := 1
		for {
			newName := fmt.Sprintf("%s_%d%s", name, counter, ext)
			newPath := BuildFilePath(user.Username, header.Path, newName)
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				filePath = newPath
				break
			}
			counter++
		}
	}

	// Save the file to disk
	outFile, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("error saving file to disk: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, *file)
	if err != nil {
		return 0, fmt.Errorf("error writing file to disk: %v", err)
	}

	// Insert file metadata into the database
	var fileID int
	query := fmt.Sprintf(`
		INSERT INTO user_files (%s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, Filename, Path, OwnerID, Size, ContentType)

	err = db.pool.QueryRow(ctx, query, header.Name, filePath, user.ID, header.Size, header.ContentType).Scan(&fileID)
	if err != nil {
		return 0, fmt.Errorf("error inserting file metadata: %v", err)
	}

	return fileID, nil
}

func (db *Database) DeleteFile(ctx context.Context, filename int) error {
	// Get the file path from the database
	var filePath string
	err := db.pool.QueryRow(ctx, `SELECT path FROM user_files WHERE filename = $1`, filename).Scan(&filePath)
	if err != nil {
		return fmt.Errorf("error retrieving file path: %v", err)
	}

	// Delete the file from disk
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error deleting file from disk: %v", err)
	}

	// Delete the file metadata from the database
	_, err = db.pool.Exec(ctx, `DELETE FROM user_files WHERE filename = $1`, filename)
	if err != nil {
		return fmt.Errorf("error deleting file metadata: %v", err)
	}

	return nil
}
