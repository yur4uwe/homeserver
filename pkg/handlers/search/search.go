package search

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/templates"
)

func FileSearchHandler(ctx context.Context, connection *db.Database, query string, directory []string, requester *user_schema.User) (*templates.SearchPageData, error) {
	search_results, err := connection.SearchFiles(ctx, query, requester, directory)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}

	return &templates.SearchPageData{
		QueryResults: true,
		Entries:      TransformEntries(search_results, nil),
	}, nil
}
