package folder

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas"
	"homeserver/pkg/schemas/folder_schema"
	"os"
)

var (
	MkdirAll  = os.MkdirAll
	RemoveAll = os.RemoveAll
	Rename    = os.Rename
)

// CREATE

func CreateFolderEntry(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error) {
	if ownerName == "" {
		return nil, fmt.Errorf("owner name cannot be empty")
	}

	if err := folder.ShouldHaveFields(ownerName)(schemas.FolderName, schemas.OwnerID, schemas.Path); err != nil {
		return nil, fmt.Errorf("invalid folder data: %v", err)
	}

	inserted_folder := folder_schema.UserFolder{}

	err := db.QueryRow(
		ctx,
		`INSERT INTO user_folders (name, parent, path, owner_id, created_at) 
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, name, parent, path, owner_id, created_at`,
		folder.Name,
		folder.Parent,
		folder.Path,
		folder.OwnerID,
	).Scan(&inserted_folder.ID, &inserted_folder.Name, &inserted_folder.Parent, &inserted_folder.Path, &inserted_folder.OwnerID, &inserted_folder.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating folder entry: %v", err)
	}

	err = MkdirAll(fmt.Sprintf("%s%s", os.Getenv(globals.EnvStoragePath), folder.Path), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating folder structure: %v", err)
	}

	return &inserted_folder, nil
}

func CreateRootFolderEntry(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error) {
	if err := folder.ShouldHaveFields(ownerName)(schemas.FolderName, schemas.OwnerID); err != nil {
		return nil, fmt.Errorf("invalid folder data: %v", err)
	}

	if folder.Path != "" {
		return nil, fmt.Errorf("root folder should not have a path")
	}
	if folder.Name != ownerName {
		return nil, fmt.Errorf("root folder name should match owner name")
	}

	inserted_folder := folder_schema.UserFolder{}

	err := db.QueryRow(
		ctx,
		`INSERT INTO user_folders (name, path, owner_id, created_at) 
		VALUES ($1, $2, $3, NOW())
		RETURNING id, name, path, owner_id, created_at`,
		folder.Name,
		folder.Path,
		folder.OwnerID,
	).Scan(&inserted_folder.ID, &inserted_folder.Name, &inserted_folder.Path, &inserted_folder.OwnerID, &inserted_folder.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating root folder entry: %v", err)
	}

	err = MkdirAll(fmt.Sprintf("%s%s", os.Getenv(globals.EnvStoragePath), folder.Path), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating root folder structure: %v", err)
	}

	return &inserted_folder, nil
}

// UPDATE

func RenameFolder(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder, newName, ownerName string) error {
	if err := folder.ShouldHaveFields(ownerName)(schemas.FolderID, schemas.OwnerID); err != nil {
		return fmt.Errorf("invalid folder data: %v", err)
	}

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
	if err := Rename(folder.AbsolutePath(), fmt.Sprintf("%s%s", os.Getenv(globals.EnvStoragePath), newPath)); err != nil {
		return fmt.Errorf("error renaming folder on disk: %v", err)
	}

	folder.Name = newName
	folder.Path = newPath

	if _, err := db.UpdateFolder(ctx, folder, ownerName); err != nil {
		return fmt.Errorf("error updating folder in database: %v", err)
	}

	return nil
}

func UpdateFolder(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error) {
	if err := folder.ShouldHaveFields(ownerName)(schemas.FolderID, schemas.FolderName, schemas.Path, schemas.FolderParent, schemas.OwnerID); err != nil {
		return nil, fmt.Errorf("invalid folder data: %v", err)
	}

	var updatedFolder folder_schema.UserFolder
	err := db.QueryRow(
		ctx,
		`UPDATE user_folders SET name = $1, parent = $2, path = $3 WHERE id = $4 AND owner_id = $5
		RETURNING id, name, parent, path, owner_id, created_at`,
		folder.Name,
		folder.Parent,
		folder.Path,
		folder.ID,
		folder.OwnerID,
	).Scan(&updatedFolder.ID, &updatedFolder.Name, &updatedFolder.Parent, &updatedFolder.Path, &updatedFolder.OwnerID, &updatedFolder.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating folder: %v", err)
	}

	return &updatedFolder, nil
}

// DELETE

func DeleteFolder(ctx context.Context, db interfaces.DatabaseInterface, folder *folder_schema.UserFolder, ownerName string) (int, error) {
	if err := folder.ShouldHaveFields(ownerName)(schemas.FolderName, schemas.OwnerID, schemas.Path); err != nil {
		return 0, fmt.Errorf("invalid folder data: %v", err)
	}

	var deletedID int
	err := db.QueryRow(
		ctx,
		`DELETE FROM user_folders WHERE name = $1 AND path = $2 AND owner_id = $3 RETURNING id`,
		folder.Name,
		folder.Path,
		folder.OwnerID,
	).Scan(&deletedID)
	if err != nil {
		return 0, fmt.Errorf("error deleting folder entry: %v", err)
	}

	if err := RemoveAll(folder.AbsolutePath()); err != nil {
		return 0, fmt.Errorf("error removing folder %s from disk: %v", folder.Name, err)
	}

	return deletedID, nil
}
