package db

import (
	"context"
	"homeserver/pkg/db/folder"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
)

func (db *Database) CreateFolderEntry(ctx context.Context, folder_data *folder_schema.UserFolder, ownerName string) error {
	return folder.CreateFolderEntry(ctx, db, folder_data, ownerName)
}

func (db *Database) CreateRootFolderEntry(ctx context.Context, folder_data *folder_schema.UserFolder, ownerName string) error {
	return folder.CreateRootFolderEntry(ctx, db, folder_data, ownerName)
}

func (db *Database) ExtractUserFolders(ctx context.Context, userID int) ([]*folder_schema.UserFolder, error) {
	return folder.ExtractUserFolders(ctx, db, userID)
}

func (db *Database) QueryFolderByName(ctx context.Context, folderName string, ownerID int, pathToFolder string) (int, error) {
	return folder.QueryFolderByName(ctx, db, folderName, ownerID, pathToFolder)
}

func (db *Database) FolderExists(ctx context.Context, folderName string, ownerID int, pathToFolder string) (bool, error) {
	return folder.FolderExists(ctx, db, folderName, ownerID, pathToFolder)
}

func (db *Database) ExtractFoldersFromDirectory(ctx context.Context, directory []string, requester *user_schema.User) ([]*folder_schema.UserFolder, error) {
	return folder.ExtractFoldersFromDirectory(ctx, db, directory, requester)
}

func (db *Database) DeleteFolder(ctx context.Context, folder_data *folder_schema.UserFolder, ownerName string) error {
	return folder.DeleteFolder(ctx, db, folder_data, ownerName)
}

func (db *Database) UpdateFolder(ctx context.Context, folder_data *folder_schema.UserFolder, ownerName string) (int, error) {
	return folder.UpdateFolder(ctx, db, folder_data, ownerName)
}

func (db *Database) RenameFolder(ctx context.Context, folder_data *folder_schema.UserFolder, newName, ownerName string) error {
	return folder.RenameFolder(ctx, db, folder_data, newName, ownerName)
}
