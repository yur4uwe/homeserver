package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/templates"
	"io"
	"mime/multipart"
	"os"
)

func BuildFilePath(owner string, path string, filename string) string {
	storage_path := os.Getenv(globals.EnvStoragePath)
	if storage_path[len(storage_path)-1] != '/' {
		panic("Storage path must end with a slash")
	}

	if path != "" && (path[len(path)-1] != '/' || path[0] != '/') {
		panic("Path must start and end with a slash")
	}

	return fmt.Sprintf("%s%s%s%s", os.Getenv(globals.EnvStoragePath), owner, path, filename)
}

func (db *Database) SaveFile(
	ctx context.Context,
	header *templates.FileInfo,
	file *multipart.File,
	save_location string,
	token string,
) (int, error) {
	user, err := db.QueryToken(ctx, token)
	if err != nil {
		return 0, fmt.Errorf("error querying token: %v", err)
	}

	// Define the file path
	filePath := BuildFilePath(user.Username, save_location, header.Name)

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
	query := `
		INSERT INTO user_files (filename, path, owner_id, metadata)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	metadata := fmt.Sprintf(`{"size": %d, "content_type": "%s"}`, header.Size, header.ContentType)
	fmt.Println("Metadata:", metadata, "query:", query)
	err = db.pool.QueryRow(ctx, query, header.Name, filePath, user.ID, metadata).Scan(&fileID)
	if err != nil {
		return 0, fmt.Errorf("error inserting file metadata: %v", err)
	}

	return fileID, nil
}

func (db *Database) DeleteFile(ctx context.Context, fileID int) error {
	// Get the file path from the database
	var filePath string
	err := db.pool.QueryRow(ctx, `SELECT path FROM user_files WHERE id = $1`, fileID).Scan(&filePath)
	if err != nil {
		return fmt.Errorf("error retrieving file path: %v", err)
	}

	// Delete the file from disk
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error deleting file from disk: %v", err)
	}

	// Delete the file metadata from the database
	_, err = db.pool.Exec(ctx, `DELETE FROM user_files WHERE id = $1`, fileID)
	if err != nil {
		return fmt.Errorf("error deleting file metadata: %v", err)
	}

	return nil
}
