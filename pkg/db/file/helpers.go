package file

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas/folder_schema"
	"os"
	"strings"
)

func EnsureFolderExists(ctx context.Context, db interfaces.DatabaseInterface, relativePath []string, ownerID int) error {
	for i, segment := range relativePath {
		if segment == "" {
			continue
		}

		exists, err := db.FolderExists(ctx, segment, ownerID, strings.Join(relativePath[:i], "/"))
		if err != nil {
			return fmt.Errorf("error checking if folder exists: %v", err)
		}

		if exists {
			continue
		}

		var parentID *int = nil
		if i > 0 {
			id, err := db.QueryFolderByName(ctx, relativePath[i-1], ownerID, strings.Join(relativePath[:i-1], "/"))
			if err != nil {
				return fmt.Errorf("error querying folder by name: %v", err)
			}
			parentID = &id
		}

		new_folder := &folder_schema.UserFolder{
			Name:    segment,
			Parent:  parentID,
			Path:    strings.Join(relativePath[:i], "/"),
			OwnerID: ownerID,
		}

		if err := db.CreateFolderEntry(ctx, new_folder); err != nil {
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
	newPath := globals.BuildFilePath(relativePathStr, newName, ext)

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, nil
		}

		newName = fmt.Sprintf("%s_%d", name, counter)
		newPath = globals.BuildFilePath(relativePathStr, newName, ext)

		counter++
	}
}
