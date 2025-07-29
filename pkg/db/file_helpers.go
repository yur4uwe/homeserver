package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"homeserver/pkg/globals"
)

func EnsureFolderExists(ctx context.Context, connection *Database, relativePath []string, ownerID int) error {
	for i, segment := range relativePath {
		if segment == "" {
			continue
		}

		exists, err := connection.FolderExists(ctx, segment, ownerID, strings.Join(relativePath[:i], "/"))
		if err != nil {
			return fmt.Errorf("error checking if folder exists: %v", err)
		}

		if exists {
			continue
		}

		var parentID *int = nil
		if i > 0 {
			id, err := connection.QueryFolderByName(ctx, relativePath[i-1], ownerID, strings.Join(relativePath[:i-1], "/"))
			if err != nil {
				return fmt.Errorf("error querying folder by name: %v", err)
			}
			parentID = &id
		}

		new_folder := &UserFolder{
			Name:    segment,
			Parent:  parentID,
			Path:    strings.Join(relativePath[:i], "/"),
			OwnerID: ownerID,
		}

		if err := connection.CreateFolderEntry(ctx, new_folder); err != nil {
			return fmt.Errorf("error creating folder entry: %v", err)
		}
	}

	return nil
}

func PreventDuplicateFileName(filePath string) (string, error) {
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

	newName := name
	newPath := BuildFilePath(relativePathStr, newName, ext)

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, nil
		}

		newName = fmt.Sprintf("%s_%d", name, counter)
		newPath = BuildFilePath(relativePathStr, newName, ext)

		counter++
	}
}

func BuildFilePath(relative_path, filename, extension string) string {
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

	return fmt.Sprintf("%s%s.%s", storage_path, file_path, extension)
}
