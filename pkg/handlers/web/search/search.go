package search

import (
	"context"
	"fmt"
	"homeserver/db"
	"homeserver/internal/templates"
	"strings"
)

func FileSearchHandler(ctx context.Context, repos *db.Database, query string, directory []string, requester *db.User) (*templates.SearchPageData, error) {
	data := &templates.SearchPageData{
		Query:        query,
		Directory:    strings.Join(directory[1:], "/"),
		URIDirectory: ParseQuerySlice(directory[1:]),
	}

	search_results, err := repos.Files.FuzzySearch(ctx, query, requester)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}

	filtered := ExcludeFilesFromParents(search_results, directory)

	data.QueryResults = true
	data.Entries = TransformEntries(filtered, nil)

	return data, nil
}
