package interfaces

import (
	"context"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"
)

type FolderInterface interface {
	ExtractFoldersFromDirectory(ctx context.Context, directory []string, requester *user_schema.User) ([]*folder_schema.UserFolder, error)
	CreateRootFolderEntry(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error)
	QueryFolderByName(ctx context.Context, folderName string, ownerID int, pathToFolder string) (*folder_schema.UserFolder, error)
	CreateFolderEntry(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error)
	UpdateFolder(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (*folder_schema.UserFolder, error)
	RenameFolder(ctx context.Context, folder *folder_schema.UserFolder, newName string, ownerName string) error
	FolderExists(ctx context.Context, folderName string, ownerID int, pathToFolder string) (bool, error)
	EnsureFolderExists(ctx context.Context, relativePath []string, owner *user_schema.User) error
	DeleteFolder(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (int, error)
	ExtractUserFolders(ctx context.Context, userID int) ([]*folder_schema.UserFolder, error)
}
