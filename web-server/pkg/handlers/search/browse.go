package search

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/web-server/pkg/templates"
	"strings"
)

func BrowseHandler(ctx context.Context, repos *repositories.Repositories, directory []string, requester *user_schema.User) (*templates.SearchPageData, error) {
	user_relative_directory := directory[1:]

	data := &templates.SearchPageData{
		Directory:    strings.Join(user_relative_directory, "/"),
		URIDirectory: ParseQuerySlice(user_relative_directory),
		QueryResults: false,
		Entries:      []templates.SearchEntry{},
	}

	file_entries, err := connection.ExtractFilesFromDirectory(ctx, strings.Join(directory, "/"), requester)
	if err != nil {
		return nil, fmt.Errorf("error extracting user files: %v", err)
	}

	folder_entries, err := connection.ExtractFoldersFromDirectory(ctx, directory, requester)
	if err != nil {
		return nil, fmt.Errorf("error extracting user folders: %v", err)
	}

	relevant_entries := TransformEntries(file_entries, folder_entries)

	SortEntries(relevant_entries)

	data.Entries = relevant_entries

	return data, nil
}
