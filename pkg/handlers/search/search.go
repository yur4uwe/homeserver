package search

import (
	"context"
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/templates"
	"strings"
)

func FileSearchHandler(ctx context.Context, connection interfaces.DatabaseInterface, query string, directory []string, requester *user_schema.User) (*templates.SearchPageData, error) {
	data := &templates.SearchPageData{
		Query:        query,
		Directory:    strings.Join(directory[1:], "/"),
		URIDirectory: ParseQuerySlice(directory[1:]),
	}

	search_results, err := connection.SearchFiles(ctx, query, requester, directory)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}

	filtered := ExcludeFilesFromParents(search_results, directory)

	data.QueryResults = true
	data.Entries = TransformEntries(filtered, nil)

	return data, nil
}
