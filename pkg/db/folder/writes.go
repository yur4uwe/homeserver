package folder

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas/folder_schema"
	"os"
)

// CREATE

func CreateFolderEntry(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO user_folders (name, parent, path, owner_id, created_at) VALUES ($1, $2, $3, $4, NOW())`,
		folder.Name,
		folder.Parent,
		folder.Path,
		folder.OwnerID,
	)
	if err != nil {
		return fmt.Errorf("error creating folder entry: %v", err)
	}

	err = os.MkdirAll(fmt.Sprintf("%s%s", os.Getenv(globals.EnvStoragePath), folder.Path), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating folder structure: %v", err)
	}

	return nil
}

// UPDATE

func RenameFolder(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder, newName string) error {
	_, err := db.Exec(
		ctx,
		`UPDATE user_folders SET name = $1 WHERE id = $2 AND owner_id = $3`,
		newName,
		folder.ID,
		folder.OwnerID,
	)
	if err != nil {
		return fmt.Errorf("error renaming folder: %v", err)
	}

	newPath := fmt.Sprintf("%s/%s", folder.Path, newName)
	if err := os.Rename(folder.AbsolutePath(), fmt.Sprintf("%s%s", os.Getenv(globals.EnvStoragePath), newPath)); err != nil {
		return fmt.Errorf("error renaming folder on disk: %v", err)
	}

	folder.Name = newName
	folder.Path = newPath

	if _, err := db.UpdateFolder(ctx, folder); err != nil {
		return fmt.Errorf("error updating folder in database: %v", err)
	}

	return nil
}

func UpdateFolder(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder) (int, error) {
	if folder.Name == "" {
		return -1, fmt.Errorf("folder name cannot be empty")
	}

	if folder.OwnerID <= 0 {
		return -1, fmt.Errorf("invalid owner ID")
	}

	if folder.Parent == nil {
		panic("user root folder cannot be updated")
	}

	_, err := db.Exec(
		ctx,
		`UPDATE user_folders SET name = $1, parent = $2, path = $3 WHERE id = $4 AND owner_id = $5`,
		folder.Name,
		folder.Parent,
		folder.Path,
		folder.ID,
		folder.OwnerID,
	)
	if err != nil {
		return -1, fmt.Errorf("error updating folder: %v", err)
	}

	return folder.ID, nil
}

// DELETE

func DeleteFolder(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder) error {
	_, err := db.Exec(
		ctx,
		`DELETE FROM user_folders WHERE id = $1 AND owner_id = $2`,
		folder.ID,
		folder.OwnerID,
	)
	if err != nil {
		return fmt.Errorf("error deleting folder entry: %v", err)
	}

	if err := os.RemoveAll(folder.AbsolutePath()); err != nil {
		return fmt.Errorf("error removing folder %s from disk: %v", folder.Name, err)
	}

	return nil
}
