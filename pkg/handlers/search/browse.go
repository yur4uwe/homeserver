package search

import (
	"context"
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/templates"
	"sort"
	"strings"
)

func BrowseHandler(ctx context.Context, connection *db.Database, directory []string, requester *user_schema.User) (*templates.SearchPageData, error) {
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

	sort.Slice(relevant_entries, func(i, j int) bool {
		if relevant_entries[i].IsFolder && !relevant_entries[j].IsFolder {
			return true
		} else if !relevant_entries[i].IsFolder && relevant_entries[j].IsFolder {
			return false
		}
		return relevant_entries[i].Name < relevant_entries[j].Name
	})

	data.Entries = relevant_entries

	return data, nil
}
