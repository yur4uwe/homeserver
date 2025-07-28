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
	storagePath := os.Getenv(globals.EnvStoragePath)

	relativePathStr := strings.ReplaceAll(filePath, "\\", "/")
	relativePathStr = strings.TrimPrefix(relativePathStr, storagePath)
	relativePathStr = strings.TrimPrefix(relativePathStr, "/")
	relativePath := strings.Split(relativePathStr, "/")
	if len(relativePath) == 0 || relativePath[0] == "" {
		panic("Invalid file path: " + filePath)
	}

	base := relativePath[len(relativePath)-1]

	relativePath = relativePath[:len(relativePath)-1]
	relativePathStr = strings.Join(relativePath, "/")

	ext := strings.Split(base, ".")[1]
	name := strings.Split(base, ".")[0]
	counter := 1

	newName := fmt.Sprintf("%s.%s", name, ext)
	newPath := BuildFilePath(relativePathStr, newName)

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, nil
		}

		newName = fmt.Sprintf("%s_%d.%s", name, counter, ext)
		newPath = BuildFilePath(relativePathStr, newName)

		counter++
	}
}

func BuildFilePath(relative_path string, filename string) string {
	storage_path := os.Getenv(globals.EnvStoragePath)
	if storage_path[len(storage_path)-1] != '/' {
		panic("Storage path must end with a slash")
	}

	var file_path string
	if relative_path != "" {
		file_path = fmt.Sprintf("%s/%s", relative_path, filename)
	} else {
		file_path = filename
	}

	return fmt.Sprintf("%s%s", storage_path, file_path)
}

func (db *Database) SaveFile(
	ctx context.Context,
	header *UserFile,
	file *multipart.File,
	requester *User,
) (int, string, error) {
	filePath := BuildFilePath(header.Path, header.Filename)

	path_slice := strings.Split(header.Path, "/")
	for i, segment := range path_slice {
		if segment == "" {
			continue
		}

		exists, err := db.FolderExists(ctx, segment, requester.ID, strings.Join(path_slice[:i], "/"))
		if err != nil {
			return 0, "", fmt.Errorf("error checking if folder exists: %v", err)
		}

		if exists {
			continue
		}

		fmt.Println("Creating new folder:", segment)

		var parentID *int = nil
		if i > 0 {
			id, err := db.QueryFolderByName(ctx, path_slice[i-1], requester.ID, strings.Join(path_slice[:i-1], "/"))
			if err != nil {
				return 0, "", fmt.Errorf("error querying folder by name: %v", err)
			}
			parentID = &id
		}

		new_folder := &UserFolder{
			Name:    segment,
			Parent:  parentID,
			Path:    strings.Join(path_slice[:i], "/"),
			OwnerID: requester.ID,
		}

		if err := db.CreateFolderEntry(ctx, new_folder); err != nil {
			return 0, "", fmt.Errorf("error creating folder entry: %v", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return 0, "", fmt.Errorf("error creating directory: %v", err)
	}

	filePath, err := PreventDuplicateFileName(filePath, requester.Username)
	if err != nil {
		return 0, "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	header.Filename = filepath.Base(filePath)

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

	return fileID, header.Filename, nil
}

func (db *Database) DeleteFile(ctx context.Context, filePath string, requester *User) error {
	filename := filepath.Base(filePath)
	if filename == "." || filename == "/" {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	path := BuildFilePath("", filePath)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error deleting file from disk: %v", err)
	}

	_, err := db.pool.Exec(ctx, `DELETE FROM user_files WHERE filename = $1 AND owner_id = $2`, filename, requester.ID)
	if err != nil {
		return fmt.Errorf("error deleting file metadata: %v", err)
	}

	return nil
}

func (db *Database) RenameFile(ctx context.Context, oldPath string, newPath string, requester *User) (string, error) {
	oldFilePath := BuildFilePath("", oldPath)
	newFilePath := BuildFilePath("", newPath)

	newFilePath, err := PreventDuplicateFileName(newFilePath, requester.Username)
	if err != nil {
		return "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	newFileName := filepath.Base(newFilePath)

	if oldFilePath == newFilePath {
		return newFilePath, nil
	}

	if err := os.Rename(oldFilePath, newFilePath); err != nil {
		return "", fmt.Errorf("error renaming file on disk: %v", err)
	}

	_, err = db.pool.Exec(ctx, `UPDATE user_files SET path = $1, filename = $2 WHERE path = $3 AND owner_id = $4`, newFilePath, newFileName, oldFilePath, requester.ID)
	if err != nil {
		return "", fmt.Errorf("error updating file metadata: %v", err)
	}

	return newFilePath, nil
}
