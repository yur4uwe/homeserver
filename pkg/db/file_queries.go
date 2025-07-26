package db

import (
	"context"
	"homeserver/pkg/templates"
)

func (db *Database) SearchFiles(ctx context.Context, query string, token string) ([]templates.SearchEntry, error) {
	// This function should implement the logic to search files in the database
	// based on the provided query. It should return a slice of SearchEntry.
	// For now, we will return an empty slice and nil error.
	return []templates.SearchEntry{}, nil
}
