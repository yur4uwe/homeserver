package tables

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
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
		}

		folders, err := db.ExtractUserFolders(ctx, user.ID)
		if err != nil {
			return "", fmt.Errorf("error extracting folders for user %s: %v", user.Username, err)
		}

		for _, folder := range folders {
			if folder.Path == "" {
				continue // Skip root folders
			}

			if _, err := db.DeleteFolder(ctx, folder, user.Username); err != nil {
				return "", fmt.Errorf("error deleting folder %s for user %s: %v", folder.Name, user.Username, err)
			}
		}
	}

	return "All entries cleared successfully", nil
}

func DropTable(ctx context.Context, db interfaces.DatabaseInterface, table_name string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s CASCADE;`, table_name)
	_, err := db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error dropping table %s: %v", table_name, err)
	}

	return nil
}
