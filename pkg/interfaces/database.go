package interfaces

import (
	"context"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/folder_schema"
	"homeserver/pkg/schemas/user_schema"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DatabaseInterface interface {
	// Basic Query Methods
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)

	// Table Management
	CreateTable(ctx context.Context, table *globals.DBTable) error
	TableExists(ctx context.Context, tableName string) (bool, error)
	ClearEntries(ctx context.Context) (string, error)
	DropTable(ctx context.Context, tableName string) error

	// File Queries
	ExtractUserFiles(ctx context.Context, userID int) ([]*file_schema.UserFile, error)
	SearchFiles(ctx context.Context, query string, requester *user_schema.User, directory []string) ([]*file_schema.UserFile, error)
	CreateFileEntry(ctx context.Context, file *file_schema.UserFile, requester *user_schema.User) (int, error)
	ExtractFilesFromDirectory(ctx context.Context, directory string, requester *user_schema.User) ([]*file_schema.UserFile, error)
	ExtractAllFiles(ctx context.Context) ([]*file_schema.UserFile, error)
	GetFile(ctx context.Context, file_data *file_schema.UserFile, requester *user_schema.User) (*file_schema.UserFile, error)
	DeleteFile(ctx context.Context, file *file_schema.UserFile, userName string) (int, error)
	RenameFile(ctx context.Context, oldPath string, newPath string, requester *user_schema.User) (string, error)
	UpdateFileEntry(ctx context.Context, file *file_schema.UserFile, ownerName string) error

	// Folder Queries
	CreateFolderEntry(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) error
	ExtractUserFolders(ctx context.Context, userID int) ([]*folder_schema.UserFolder, error)
	QueryFolderByName(ctx context.Context, folderName string, ownerID int, pathToFolder string) (int, error)
	FolderExists(ctx context.Context, folderName string, ownerID int, pathToFolder string) (bool, error)
	ExtractFoldersFromDirectory(ctx context.Context, directory []string, requester *user_schema.User) ([]*folder_schema.UserFolder, error)
	DeleteFolder(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) error
	UpdateFolder(ctx context.Context, folder *folder_schema.UserFolder, ownerName string) (int, error)
	RenameFolder(ctx context.Context, folder *folder_schema.UserFolder, newName string, ownerName string) error
	EnsureFolderExists(ctx context.Context, relativePath []string, owner *user_schema.User) error

	// Token Queries
	CheckToken(ctx context.Context, token_str string) (int, error)
	CreateToken(ctx context.Context, userID int, token_str string) error
	QueryToken(ctx context.Context, token_str string) (*user_schema.User, error)
	VerifyToken(ctx context.Context, token_str string) error
	DeleteToken(ctx context.Context, token_str string) error

	// User Queries
	CreateUserEntry(ctx context.Context, user *user_schema.User) (*user_schema.User, error)
	ExtractUserEntries(ctx context.Context) ([]*user_schema.User, error)
	CheckUserFolder(ctx context.Context, user *user_schema.User) error
	QueryUserOnLogin(ctx context.Context, form globals.LoginForm) (int, bool, error)
	QueryUsername(ctx context.Context, username string) (int, error)
	QueryEmail(ctx context.Context, email string) (int, error)
	QueryPassword(ctx context.Context, password string) (int, error)
}
