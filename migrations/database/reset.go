package database

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
)

func ResetFolderStructure(connection *db.Database) error {
	ctx := context.Background()

	result, err := connection.ClearEntries(ctx)
	if err != nil {
		return fmt.Errorf("error clearing folder entries: %v", err)
	}

	fmt.Println(result)

	users, err := connection.ExtractUserEntries(ctx)
	if err != nil {
		return fmt.Errorf("error extracting user entries: %v", err)
	}

	for _, user := range users {
		if err := connection.CheckUserFolder(ctx, user); err != nil {
			return fmt.Errorf("error checking user folder for %s: %v", user.Username, err)
		}
	}

	return nil
}
