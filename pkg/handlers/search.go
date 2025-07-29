package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
)

func BrowseHandler(ctx context.Context, connection *db.Database, directory []string, requester *db.User) (*templates.SearchPageData, error) {
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

func FileSearchHandler(ctx context.Context, connection *db.Database, query string, directory []string, requester *db.User) (*templates.SearchPageData, error) {
	search_results, err := connection.SearchFiles(ctx, query, requester, directory)
	if err != nil {
		return nil, fmt.Errorf("error searching files: %v", err)
	}

	return &templates.SearchPageData{
		QueryResults: true,
		Entries:      TransformEntries(search_results, nil),
	}, nil
}

func SearchHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		query := r.URL.Query().Get("query")
		directory_string, err := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/search/"))
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API url.QueryUnescape error: "+err.Error(), "SearchHandler:30")
			utils.ErrorResponse(w, "Internal Server Error: Failed to unescape directory", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Search API accessed: query: '" + query + "', directory: '" + directory_string + "'")

		directory := []string{reqCont.Requester.Username}
		if directory_string != "" {
			directory = append(directory, strings.Split(directory_string, "/")...)
		}

		var data *templates.SearchPageData = nil

		if len(query) > 0 {
			data, err = FileSearchHandler(r.Context(), connection, query, directory, reqCont.Requester)
		} else {
			data, err = BrowseHandler(r.Context(), connection, directory, reqCont.Requester)
		}

		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			utils.ErrorResponse(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
			return
		}

		var parent []string
		if len(directory[1:]) > 0 {
			parent = directory[1 : len(directory)-1]
		}

		data.Parent = strings.Join(parent, "/")
		data.URIParent = ParseQuerySlice(parent)

		fmt.Println("SearchHandler: data:", data)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		var tmpl *template.Template
		if reqCont.ToRestorePage || reqCont.RerenderBody || r.URL.Query().Get("redirect") != "" {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_base.html", "../static/templates/search_results.html"))

			w.Header().Set("HX-Push-Url", "/search/"+url.QueryEscape(strings.Join(directory[1:], "/")))
			tmpl.ExecuteTemplate(w, "search_base.html", &data)
		} else {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_results.html"))

			tmpl.Execute(w, data)
		}
	}
}
