package search

import (
	"context"
	"fmt"
	"homeserver/db"
	"homeserver/internal/templates"
	"net/url"
	"strconv"
	"strings"
)

const pageSize = 20

func FileSearchHandler(ctx context.Context, repos *db.Database, query string, directory []string, requester *db.User, page int32) (*templates.SearchPageData, error) {
	data := &templates.SearchPageData{
		Query:        query,
		Directory:    strings.Join(directory[1:], "/"),
		URIDirectory: ParseQuerySlice(directory[1:]),
	}

	offset := (int64(page) - 1) * pageSize
	files, err := repos.Files.FuzzySearch(ctx, query, strings.Join(directory, "/"), requester, pageSize+1, offset)
	if err != nil {
		return nil, fmt.Errorf("search files: %w", err)
	}

	path := "/search/" + data.URIDirectory
	pageURL := func(number int64) string {
		return path + "?" + url.Values{"query": {query}, "page": {strconv.FormatInt(number, 10)}}.Encode()
	}
	if page > 1 {
		data.PreviousPageURL = pageURL(int64(page) - 1)
	}
	if len(files) > pageSize {
		data.NextPageURL = pageURL(int64(page) + 1)
		files = files[:pageSize]
	}

	data.QueryResults = true
	data.Entries = TransformEntries(files, nil)

	return data, nil
}
