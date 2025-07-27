package db

import (
	"context"
	"fmt"
	"strings"
)

func (db *Database) ExtractUserFiles(ctx context.Context, userID int) ([]*UserFile, error) {
	rows, err := db.pool.Query(
		ctx,
		fmt.Sprintf(`SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s FROM user_files WHERE %s = $1`,
			FileID, Filename, Path, Size, ContentType, OwnerID, CreatedAt, UpdatedAt, Metadata, OwnerID,
		),
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying user files: %v", err)
	}
	defer rows.Close()

	var entries []*UserFile
	for rows.Next() {
		entry := &UserFile{}
		if err := rows.Scan(&entry.ID, &entry.Filename, &entry.Path, &entry.Size, &entry.ContentType, &entry.OwnerID, &entry.CreatedAt, &entry.UpdatedAt, &entry.Metadata); err != nil {
			return nil, fmt.Errorf("error scanning user file entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (db *Database) SearchFiles(ctx context.Context, query string, token string) ([]*UserFile, error) {
	user, err := db.QueryToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("error querying user by token: %v", err)
	}

	files, err := db.ExtractUserFiles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error extracting user files: %v", err)
	}

	// Example static data for demonstration: only files, no folders
	// results := []templates.SearchEntry{
	// 	{
	// 		Name:     "example_file.txt",
	// 		Path:     "/files/example_file.txt",
	// 		IsFolder: false,
	// 		Size:     "1 MB",
	// 	},
	// 	{
	// 		Name:     "photo.jpg",
	// 		Path:     "/files/photos/photo.jpg",
	// 		IsFolder: false,
	// 		Size:     "54 KB",
	// 	},
	// 	{
	// 		Name:     "document.pdf",
	// 		Path:     "/files/docs/document.pdf",
	// 		IsFolder: false,
	// 		Size:     "2 MB",
	// 	},
	// 	{
	// 		Name:     "music.mp3",
	// 		Path:     "/files/music/music.mp3",
	// 		IsFolder: false,
	// 		Size:     "4 MB",
	// 	},
	// 	{
	// 		Name:     "presentation.pptx",
	// 		Path:     "/files/presentations/presentation.pptx",
	// 		IsFolder: false,
	// 		Size:     "1 GB",
	// 	},
	// }
	return files, nil
}

func (db *Database) CreateFileEntry(ctx context.Context, file *UserFile, requester *User) error {
	if file.Filename == "" || file.Size == 0 || file.ContentType == "" || requester == nil {
		return fmt.Errorf("invalid file entry data")
	}

	// Ensure the file path is valid
	if strings.Contains(file.Path, "..") || strings.Contains(file.Path, "//") || strings.Contains(file.Path, "\\") {
		return fmt.Errorf("invalid file path: %s", file.Path)
	}

	relative_path := strings.Split(file.Path, "/")
	if len(relative_path) == 0 || relative_path[0] == "" {
		relative_path = []string{requester.Username} // Ensure at least one element
	}

	id, err := db.QueryFolderByName(ctx, relative_path[len(relative_path)-1], requester.ID, strings.Join(relative_path[:len(relative_path)-1], "/"))
	if err != nil {
		return fmt.Errorf("error querying folder by name: %v", err)
	}
	file.Directory = id

	var fileID int
	query := fmt.Sprintf(`
		INSERT INTO user_files (%s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, Filename, Path, OwnerID, Size, ContentType)

	err = db.pool.QueryRow(ctx, query, file.Filename, strings.Join(relative_path, "/"), requester.ID, file.Size, file.ContentType).Scan(&fileID)
	if err != nil {
		return fmt.Errorf("error inserting file metadata: %v", err)
	}

	return nil
}
