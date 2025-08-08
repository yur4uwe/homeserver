package file

import (
	"context"
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/utils"
	"strings"

	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	MkdirAll  = os.MkdirAll
	Rename    = os.Rename
	RemoveAll = os.RemoveAll
	Create    = os.Create
	Remove    = os.Remove
	Copy      = io.Copy
)

// CREATE

func SaveFile(
	ctx context.Context,
	db interfaces.DatabaseInterface,
	header *file_schema.UserFile,
	file *multipart.File,
	requester *user_schema.User,
) (int, string, error) {
	if err := requester.ShouldHaveFields(schemas.UserID); err != nil {
		return 0, "", fmt.Errorf("requester must have a valid user ID: %v", err)
	}

	if err := header.ShouldHaveFields(requester.Username)(schemas.Filename, schemas.Extension, schemas.Path); err != nil {
		return 0, "", fmt.Errorf("file header must have valid fields: %v", err)
	}

	filePath := globals.BuildFilePath(header.Path, header.Filename, header.Extension)

	path_slice := strings.Split(header.Path, "/")
	if err := db.EnsureFolderExists(ctx, path_slice, requester); err != nil {
		return 0, "", fmt.Errorf("error ensuring parent folders exist: %v", err)
	}

	if err := MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return 0, "", fmt.Errorf("error creating directory: %v", err)
	}

	filePath, err := PreventDuplicateFileName(filePath)
	if err != nil {
		return 0, "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	header.Filename = strings.Split(filepath.Base(filePath), ".")[0]

	outFile, err := Create(filePath)
	if err != nil {
		return 0, "", fmt.Errorf("error saving file to disk: %v", err)
	}
	defer outFile.Close()

	_, err = Copy(outFile, *file)
	if err != nil {
		return 0, "", fmt.Errorf("error writing file to disk: %v", err)
	}

	fileID, err := db.CreateFileEntry(ctx, header, requester)
	if err != nil {
		return 0, "", fmt.Errorf("error creating file entry: %v", err)
	}

	return fileID, header.Filename + "." + header.Extension, nil
}

func CreateFileEntry(ctx context.Context, db interfaces.DatabaseInterface, file *file_schema.UserFile, requester *user_schema.User) (int, error) {
	if err := file.ShouldHaveFields(requester.Username)(schemas.Filename, schemas.Extension, schemas.Path, schemas.OwnerID, schemas.Size, schemas.ContentType); err != nil {
		return -1, fmt.Errorf("invalid file data: %v", err)
	}

	relativePath := strings.Split(file.Path, "/")
	if len(relativePath) == 0 || relativePath[0] == "" {
		relativePath = []string{requester.Username}
	}

	directory, err := db.QueryFolderByName(ctx, relativePath[len(relativePath)-1], requester.ID, strings.Join(relativePath[:len(relativePath)-1], "/"))
	if err != nil {
		return -1, fmt.Errorf("error querying folder by name: %v", err)
	}
	file.Directory = directory.ID

	var fileID int
	query := fmt.Sprintf(`
		INSERT INTO user_files (%s, %s, %s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, schemas.Filename, schemas.Extension, schemas.Path, schemas.OwnerID, schemas.Size, schemas.ContentType, schemas.Directory)

	err = db.QueryRow(ctx, query, file.Filename, file.Extension, strings.Join(relativePath, "/"), requester.ID, file.Size, file.ContentType, file.Directory).Scan(&fileID)
	if err != nil {
		return -1, fmt.Errorf("error inserting file metadata: %v", err)
	}

	return fileID, nil
}

// UPDATE

func RenameFile(ctx context.Context, db interfaces.DatabaseInterface, oldPath string, newName string, requester *user_schema.User) (string, error) {
	if oldPath == "" || newName == "" {
		return "", fmt.Errorf("invalid file path: oldPath and newName cannot be empty")
	}

	err := requester.ShouldHaveFields(schemas.Username)
	if err != nil {
		return "", fmt.Errorf("error validating requester: %v", err)
	}

	oldRelFilePath, err := utils.ParsePathWithUName(oldPath, requester.Username)
	if err != nil {
		return "", fmt.Errorf("error parsing old file path: %v", err)
	}

	newRelFilePath := append(oldRelFilePath[:len(oldRelFilePath)-1], newName)

	oldAbsFilePath := globals.BuildFilePath(requester.Username, strings.Split(oldPath, ".")[0], strings.Split(oldPath, ".")[1])
	newAbsFilePath, err := PreventDuplicateFileName(strings.Join(newRelFilePath, "/"))
	if err != nil {
		return "", fmt.Errorf("error preventing duplicate file name: %v", err)
	}

	newFileName := strings.Split(filepath.Base(newAbsFilePath), ".")[0]
	oldFileName := strings.Split(filepath.Base(oldAbsFilePath), ".")[0]

	if strings.Join(oldRelFilePath, "/") == strings.Join(newRelFilePath, "/") {
		return strings.Join(newRelFilePath, "/"), nil
	}

	if err := Rename(oldAbsFilePath, newAbsFilePath); err != nil {
		return "", fmt.Errorf("error renaming file on disk: %v", err)
	}

	relPath := strings.Join(newRelFilePath[:len(newRelFilePath)-1], "/")

	tag, err := db.Exec(
		ctx,
		`UPDATE user_files SET path = $1, filename = $2, updated_at = NOW() WHERE filename = $3 AND path = $4 AND owner_id = $5`,
		relPath, newFileName, oldFileName, relPath, requester.ID,
	)
	if err != nil {
		return "", fmt.Errorf("error updating file metadata: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return "", fmt.Errorf("no file found to rename: %s", oldRelFilePath)
	}

	return newRelFilePath[len(newRelFilePath)-1], nil
}

func UpdateFileEntry(ctx context.Context, db interfaces.DatabaseInterface, file *file_schema.UserFile, ownerName string) error {
	if err := file.ShouldHaveFields(ownerName)(schemas.FileID, schemas.Filename, schemas.Extension, schemas.Path, schemas.Size, schemas.ContentType); err != nil {
		return fmt.Errorf("invalid file data: %v", err)
	}

	query := `
		UPDATE user_files
		SET filename = $1, extension = $2, path = $3, size = $4, content_type = $5, updated_at = NOW()
		WHERE id = $6
	`

	tag, err := db.Exec(ctx, query, file.Filename, file.Extension, file.Path, file.Size, file.ContentType, file.ID)
	if err != nil {
		return fmt.Errorf("error updating file entry: %v", err)
	}

	if tag.RowsAffected() < 1 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

// DELETE

func DeleteFile(ctx context.Context, db interfaces.DatabaseInterface, file *file_schema.UserFile, user_name string) (int, error) {
	if user_name == "" {
		return 0, fmt.Errorf("user name cannot be empty")
	}

	if err := file.ShouldHaveFields(user_name)(schemas.Filename, schemas.OwnerID, schemas.Path); err != nil {
		return 0, fmt.Errorf("invalid file data: %v", err)
	}

	path := file.AbsolutePath()

	if err := Remove(path); err != nil {
		return 0, fmt.Errorf("error deleting file from disk: %v", err)
	}

	var id int
	err := db.QueryRow(
		ctx,
		`DELETE FROM user_files WHERE filename = $1 AND owner_id = $2 RETURNING id`,
		file.Filename, file.OwnerID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error deleting file metadata: %v", err)
	}

	return id, nil
}
