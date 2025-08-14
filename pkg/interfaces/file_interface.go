package interfaces

import (
	"context"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"mime/multipart"
)

type FileInterface interface {
	ExtractUserFiles(ctx context.Context, userID int) ([]*file_schema.UserFile, error)
	SearchFiles(ctx context.Context, query string, requester *user_schema.User, directory []string) ([]*file_schema.UserFile, error)
	CreateFileEntry(ctx context.Context, file *file_schema.UserFile, requester *user_schema.User) (int, error)
	ExtractFilesFromDirectory(ctx context.Context, directory string, requester *user_schema.User) ([]*file_schema.UserFile, error)
	ExtractAllFiles(ctx context.Context) ([]*file_schema.UserFile, error)
	GetFile(ctx context.Context, file_data *file_schema.UserFile, requester *user_schema.User) (*file_schema.UserFile, error)
	DeleteFile(ctx context.Context, file *file_schema.UserFile, userName string) (int, error)
	RenameFile(ctx context.Context, oldFile *file_schema.UserFile, newPath string, requester *user_schema.User) (string, error)
	UpdateFileEntry(ctx context.Context, file *file_schema.UserFile, ownerName string) error
	SaveFile(ctx context.Context, header *file_schema.UserFile, fileData *multipart.File, requester *user_schema.User) (int, string, error)
}
