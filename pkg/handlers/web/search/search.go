package search

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/internal/templates"
	"strings"
)

func FileSearchHandler(ctx context.Context, repos *repositories.Repositories, query string, directory []string, requester *uschema.User) (*templates.SearchPageData, error) {
	data := &templates.SearchPageData{
		Query:        query,
		Directory:    strings.Join(directory[1:], "/"),
		URIDirectory: ParseQuerySlice(directory[1:]),
	}

	search_results, err := repos.FileRepo.FuzzySearch(ctx, query, strings.Join(directory, "/"), requester)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}

	filtered := ExcludeFilesFromParents(search_results, directory)

	data.QueryResults = true
	data.Entries = TransformEntries(filtered, nil)

	return data, nil
}
