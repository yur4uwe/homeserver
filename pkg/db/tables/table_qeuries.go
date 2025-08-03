package tables

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
	"os"
)

func ClearEntries(ctx context.Context, db interfaces.DatabaseInterface) (string, error) {
	users, err := db.ExtractUserEntries(ctx)
	if err != nil {
		return "", fmt.Errorf("error extracting user entries: %v", err)
	}

	for _, user := range users {
		files, err := db.ExtractUserFiles(ctx, user.ID)
		if err != nil {
			return "", fmt.Errorf("error extracting files for user %s: %v", user.Username, err)
		}

		for _, file := range files {
			if _, err := db.DeleteFile(ctx, file, user.Username); err != nil {
				return "", fmt.Errorf("error deleting file %s for user %s: %v", file.Filename, user.Username, err)
			}

			if err := os.Remove(file.AbsolutePath()); err != nil {
				return "", fmt.Errorf("error removing file %s from disk: %v", file.Filename, err)
			}
		}

		folders, err := db.ExtractUserFolders(ctx, user.ID)
		if err != nil {
			return "", fmt.Errorf("error extracting folders for user %s: %v", user.Username, err)
		}

		for _, folder := range folders {
			if err := db.DeleteFolder(ctx, folder); err != nil {
				return "", fmt.Errorf("error deleting folder %s for user %s: %v", folder.Name, user.Username, err)
			}

			if err := os.RemoveAll(folder.AbsolutePath()); err != nil {
				return "", fmt.Errorf("error removing folder %s from disk: %v", folder.Name, err)
			}
		}
	}

	return "All entries cleared successfully", nil
}
