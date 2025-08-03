package folder

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"strings"
)

func ExtractUserFolders(ctx context.Context, db interfaces.DatabaseInterface, userID int) ([]*folder_schema.UserFolder, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, name, parent, owner_id, created_at FROM user_folders WHERE owner_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying user folders: %v", err)
	}
	defer rows.Close()

	var entries []*folder_schema.UserFolder
	for rows.Next() {
		entry := &folder_schema.UserFolder{}
		if err := rows.Scan(&entry.ID, &entry.Name, &entry.Parent, &entry.OwnerID, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning user folder entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func QueryFolderByName(ctx context.Context, db interfaces.DatabaseInterface, folderName string, ownerID int, pathToFolder string) (int, error) {
	var folderID int
	err := db.QueryRow(
		ctx,
		`SELECT id FROM user_folders WHERE name = $1 AND owner_id = $2 AND path = $3`,
		folderName,
		ownerID,
		pathToFolder,
	).Scan(&folderID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return -1, fmt.Errorf("folder does not exist: %s", folderName)
		}
		return -1, fmt.Errorf("error querying folder by name: %v", err)
	}
	return folderID, nil
}

func FolderExists(ctx context.Context, db interfaces.DatabaseInterface, folderName string, ownerID int, pathToFolder string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM user_folders WHERE name = $1 AND owner_id = $2 AND path = $3)`,
		folderName,
		ownerID,
		pathToFolder,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if folder exists: %v", err)
	}
	return exists, nil
}

func ExtractFoldersFromDirectory(ctx context.Context, db interfaces.DatabaseInterface, directory []string, requester *user_schema.User) ([]*folder_schema.UserFolder, error) {
	query := fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s
		FROM user_folders
		WHERE path = $1 AND owner_id = $2
	`, schemas.FolderID, schemas.FolderName, schemas.Path, schemas.FolderParent, schemas.OwnerID, schemas.CreatedAt)

	rows, err := db.Query(ctx, query, strings.Join(directory, "/"), requester.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying folders in directory: %v", err)
	}
	defer rows.Close()

	var folders []*folder_schema.UserFolder
	for rows.Next() {
		folder := &folder_schema.UserFolder{}
		if err := rows.Scan(&folder.ID, &folder.Name, &folder.Path, &folder.Parent, &folder.OwnerID, &folder.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning folder entry: %v", err)
		}
		folders = append(folders, folder)
	}

	return folders, nil
}
