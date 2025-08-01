package db

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
	"os"
	"strings"
)

// CreateFolderEntry creates a new folder entry in the database and ensures the folder structure exists on disk.
// WRITES TO DISK
func (db *Database) CreateFolderEntry(ctx context.Context, folder *folder_schema.UserFolder) error {
	_, err := db.pool.Exec(
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

func (db *Database) ExtractUserFolders(ctx context.Context, userID int) ([]*folder_schema.UserFolder, error) {
	rows, err := db.pool.Query(
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

func (db *Database) QueryFolderByName(ctx context.Context, folderName string, ownerID int, path_to_folder string) (int, error) {
	var folderID int
	err := db.pool.QueryRow(
		ctx,
		`SELECT id FROM user_folders WHERE name = $1 AND owner_id = $2 AND path = $3`,
		folderName,
		ownerID,
		path_to_folder,
	).Scan(&folderID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return -1, fmt.Errorf("folder does not exist: %s", folderName)
		}
		return -1, fmt.Errorf("error querying folder by name: %v", err)
	}
	return folderID, nil
}

func (db *Database) FolderExists(ctx context.Context, folderName string, ownerID int, path_to_folder string) (bool, error) {
	var exists bool
	err := db.pool.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM user_folders WHERE name = $1 AND owner_id = $2 AND path = $3)`,
		folderName,
		ownerID,
		path_to_folder,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if folder exists: %v", err)
	}
	return exists, nil
}

func (db *Database) ExtractFoldersFromDirectory(ctx context.Context, directory []string, requester *user_schema.User) ([]*folder_schema.UserFolder, error) {
	query := fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s, %s
		FROM user_folders
		WHERE path = $1 AND owner_id = $2
	`, schemas.FolderID, schemas.FolderName, schemas.Path, schemas.FolderParent, schemas.OwnerID, schemas.CreatedAt)

	rows, err := db.pool.Query(ctx, query, strings.Join(directory, "/"), requester.ID)
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

func (db *Database) DeleteFolder(ctx context.Context, folder *folder_schema.UserFolder) error {
	db.pool.Exec(
		ctx,
		`DELETE FROM user_folders WHERE id = $1 AND owner_id = $2`,
		folder.ID,
		folder.OwnerID,
	)

	if err := os.RemoveAll(folder.AbsolutePath()); err != nil {
		return fmt.Errorf("error removing folder %s from disk: %v", folder.Name, err)
	}

	return nil
}
