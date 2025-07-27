package main

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"os"
	"strings"
)

func CreateFolderStructure(connection *db.Database) error {
	ctx := context.Background()

	users, err := connection.ExtractUserEntries(ctx)
	if err != nil {
		return fmt.Errorf("error extracting user entries: %v", err)
	}

	for _, user := range users {
		files, err := connection.ExtractUserFiles(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("error extracting files for user %s: %v", user.Username, err)
		}

		for _, file := range files {
			path_string := strings.ReplaceAll(file.Path, "\\", "/")
			relative_path_string := strings.TrimPrefix(path_string, os.Getenv(globals.EnvStoragePath))
			relative_path := strings.Split(relative_path_string, "/")
			path_to_file := relative_path[:len(relative_path)-1]

			for i, folder_name := range path_to_file {
				var parentID *int = nil
				if i > 0 {
					parent_name := path_to_file[i-1]
					parent_path := strings.Join(path_to_file[:i-1], "/")

					id, err := connection.QueryFolderByName(ctx, parent_name, user.ID, parent_path)
					if err != nil {
						return fmt.Errorf("error querying parent folder %s for user %s: %v", parent_name, user.Username, err)
					}
					parentID = &id
				}

				folder := &db.UserFolder{
					Name:    folder_name,
					Parent:  parentID,
					Path:    strings.Join(path_to_file[:i], "/"),
					OwnerID: user.ID,
				}

				_, err := connection.QueryFolderByName(ctx, folder.Name, user.ID, folder.Path)
				if err != nil {
					if strings.Contains(err.Error(), "folder does not exist") {
						if err := connection.CreateFolderEntry(ctx, folder); err != nil {
							return fmt.Errorf("error creating folder %s for user %s: %v", folder.Name, user.Username, err)
						}
					} else {
						return fmt.Errorf("error querying folder by name: %v", err)
					}
				}
			}

		}
	}

	return nil
}
