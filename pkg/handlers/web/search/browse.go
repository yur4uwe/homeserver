package search

import (
	"context"
	"fmt"
	"homeserver/db/repositories"
	uschema "homeserver/db/repositories/user_repo/schema"
	"homeserver/internal/templates"
	"strings"
)

func BrowseHandler(ctx context.Context, repos *repositories.Repositories, directory []string, requester *uschema.User) (*templates.SearchPageData, error) {
	user_relative_directory := directory[1:]

	data := &templates.SearchPageData{
		Directory:    strings.Join(user_relative_directory, "/"),
		URIDirectory: ParseQuerySlice(user_relative_directory),
		QueryResults: false,
		Entries:      []templates.SearchEntry{},
	}

	file_entries, err := repos.FileRepo.GetFilesByDirectory(ctx, strings.Join(directory, "/"), requester)
	if err != nil {
		return nil, fmt.Errorf("error extracting user files: %v", err)
	}

	folder_entries, err := repos.FolderRepo.GetChildren(
		ctx,
		strings.Join(directory[:len(directory)-1], "/"),
		directory[len(directory)-1],
		requester,
	)
	if err != nil {
		return nil, fmt.Errorf("error extracting user folders: %v", err)
	}

	relevant_entries := TransformEntries(file_entries, folder_entries)

	SortEntries(relevant_entries)

	data.Entries = relevant_entries

	return data, nil
}
