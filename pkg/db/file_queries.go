package db

import (
	"context"
	"fmt"
)

func (db *Database) ExtractUserFiles(ctx context.Context, userID int) ([]*UserFiles, error) {
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

	var entries []*UserFiles
	for rows.Next() {
		entry := &UserFiles{}
		if err := rows.Scan(&entry.ID, &entry.Filename, &entry.Path, &entry.Size, &entry.ContentType, &entry.OwnerID, &entry.CreatedAt, &entry.UpdatedAt, &entry.Metadata); err != nil {
			return nil, fmt.Errorf("error scanning user file entry: %v", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (db *Database) SearchFiles(ctx context.Context, query string, token string) ([]*UserFiles, error) {
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
