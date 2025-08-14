package db

import (
	"context"
	"homeserver/pkg/db/file"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"mime/multipart"
)

func (db *Database) ExtractUserFiles(ctx context.Context, userID int) ([]*file_schema.UserFile, error) {
	return file.ExtractUserFiles(ctx, db, userID)
}

func (db *Database) SearchFiles(ctx context.Context, query string, requester *user_schema.User, directory []string) ([]*file_schema.UserFile, error) {
	return file.SearchFiles(ctx, db, query, requester, directory)
}

func (db *Database) CreateFileEntry(ctx context.Context, file_data *file_schema.UserFile, requester *user_schema.User) (int, error) {
	return file.CreateFileEntry(ctx, db, file_data, requester)
}

func (db *Database) ExtractFilesFromDirectory(ctx context.Context, directory string, requester *user_schema.User) ([]*file_schema.UserFile, error) {
	return file.ExtractFilesFromDirectory(ctx, db, directory, requester)
}

func (db *Database) ExtractAllFiles(ctx context.Context) ([]*file_schema.UserFile, error) {
	return file.ExtractAllFiles(ctx, db)
}

func (db *Database) GetFile(ctx context.Context, file_data *file_schema.UserFile, requester *user_schema.User) (*file_schema.UserFile, error) {
	return file.GetFile(ctx, db, file_data, requester)
}

func (db *Database) DeleteFile(ctx context.Context, file_data *file_schema.UserFile, userName string) (int, error) {
	return file.DeleteFile(ctx, db, file_data, userName)
}

// oldPath is the current path of the file, newPath is the desired new path.
// It returns the new path on success.
// oldPath and newPath should include the filename and extension but not the user directory.
// For example, if the file is located at "testuser/testfile.txt", oldPath
// should be "testfile.txt" and newPath should be "renamedfile.txt".
func (db *Database) RenameFile(ctx context.Context, oldFile *file_schema.UserFile, newPath string, requester *user_schema.User) (string, error) {
	return file.RenameFile(ctx, db, oldFile, newPath, requester)
}

func (db *Database) UpdateFileEntry(ctx context.Context, file_data *file_schema.UserFile, ownerName string) error {
	return file.UpdateFileEntry(ctx, db, file_data, ownerName)
}

func (db *Database) SaveFile(ctx context.Context, header *file_schema.UserFile, fileData *multipart.File, requester *user_schema.User) (int, string, error) {
	return file.SaveFile(ctx, db, header, fileData, requester)
}

func (db *Database) EnsureFolderExists(ctx context.Context, relativePath []string, owner *user_schema.User) error {
	return file.EnsureFolderExists(ctx, db, relativePath, owner)
}
