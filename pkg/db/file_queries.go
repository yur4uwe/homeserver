package db

import (
	"context"
	"homeserver/pkg/templates"
)

func (db *Database) SearchFiles(ctx context.Context, query string, token string) ([]templates.SearchEntry, error) {
	// Example static data for demonstration: only files, no folders
	results := []templates.SearchEntry{
		{
			Name:     "example_file.txt",
			Path:     "/files/example_file.txt",
			IsFolder: false,
			Size:     "1 MB",
		},
		{
			Name:     "photo.jpg",
			Path:     "/files/photos/photo.jpg",
			IsFolder: false,
			Size:     "54 KB",
		},
		{
			Name:     "document.pdf",
			Path:     "/files/docs/document.pdf",
			IsFolder: false,
			Size:     "2 MB",
		},
		{
			Name:     "music.mp3",
			Path:     "/files/music/music.mp3",
			IsFolder: false,
			Size:     "4 MB",
		},
		{
			Name:     "presentation.pptx",
			Path:     "/files/presentations/presentation.pptx",
			IsFolder: false,
			Size:     "1 GB",
		},
	}
	return results, nil
}
