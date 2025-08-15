package database

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"strings"
)

func SeparateExtensionsFromNames(connection *db.Database) error {
	ctx := context.Background()

	users, err := connection.ExtractUserEntries(ctx)
	if err != nil {
		return fmt.Errorf("error extracting user entries: %v", err)
	}

	for _, user := range users {
		files, err := connection.ExtractUserFiles(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("error extracting user files: %v", err)
		}

		for _, file := range files {

			user_sig := strings.Split(user.Username, ".")
			if len(user_sig) > 1 && file.Extension == "" {
				user.Username = strings.Join(user_sig[:len(user_sig)-1], ".")
				file.Extension = user_sig[len(user_sig)-1]
			}

			if err := connection.UpdateFileEntry(ctx, file, user.Username); err != nil {
				return fmt.Errorf("error updating file entry: %v", err)
			}
		}
	}

	return nil
}
