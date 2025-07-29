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

		reqCont.LogHistory.Log("[Handler] Search API accessed")

		query := r.URL.Query().Get("query")
		directory_string, err := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/search/"))
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API url.QueryUnescape error: "+err.Error(), "SearchHandler:30")
			utils.ErrorResponse(w, "Internal Server Error: Failed to unescape directory", http.StatusInternalServerError)
			return
		}

		directory := []string{reqCont.Requester.Username}
		if directory_string != "" {
			directory = append(directory, strings.Split(directory_string, "/")...)
		}

		// remote_path := strings.TrimPrefix(r.Header.Get("Referer"), fmt.Sprintf("http://%s/search/", r.Host))

		// directory_without_username := directory[1:]

		// if remote_path != strings.Join(directory_without_username, "/") {
		// 	reqCont.LogHistory.Log(fmt.Sprintf("[Handler] Restoring directory for client: '%s'", directory))
		// 	if len(directory) > 0 {
		// 		directory_string = ParseQuerySlice(directory_without_username)
		// 	}

		// 	w.Header().Set("HX-Push-Url", "/search/"+directory_string)
		// }

		// reqCont.LogHistory.Log(fmt.Sprintf("[Handler] Search API query: '%s', directory: '%s'", query, strings.Join(directory, "/")))

		// file_entries, err := connection.SearchFiles(r.Context(), query, reqCont.Requester, directory)
		// if err != nil {
		// 	reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
		// 	utils.ErrorResponse(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
		// 	return
		// }

		// folder_entries, err := connection.ExtractFoldersFromDirectory(r.Context(), directory, reqCont.Requester)
		// if err != nil {
		// 	reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
		// 	utils.ErrorResponse(w, "Internal Server Error: Failed to search folders", http.StatusInternalServerError)
		// 	return
		// }

		// relevant_entries := TransformEntries(file_entries, folder_entries)

		// sort.Slice(relevant_entries, func(i, j int) bool {
		// 	if relevant_entries[i].IsFolder && !relevant_entries[j].IsFolder {
		// 		return true
		// 	} else if !relevant_entries[i].IsFolder && relevant_entries[j].IsFolder {
		// 		return false
		// 	}
		// 	return relevant_entries[i].Name < relevant_entries[j].Name
		// })

		// var parent []string
		// if len(directory_without_username) > 0 {
		// 	parent = directory_without_username[:len(directory_without_username)-1]
		// }

		// data := &templates.SearchPageData{
		// 	Query:        query,
		// 	Parent:       strings.Join(parent, "/"),
		// 	Directory:    strings.Join(directory_without_username, "/"),
		// 	URIQuery:     url.QueryEscape(query),
		// 	URIParent:    ParseQuerySlice(parent),
		// 	URIDirectory: ParseQuerySlice(directory_without_username),
		// 	QueryResults: len(query) > 0,
		// 	Entries:      relevant_entries,
		// }

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

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		var tmpl *template.Template
		if reqCont.ToRestorePage || reqCont.RerenderBody || r.URL.Query().Get("redirect") != "" {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_base.html", "../static/templates/search_results.html"))

			tmpl.ExecuteTemplate(w, "search_base.html", data)
		} else {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_results.html"))

			tmpl.Execute(w, data)
		}
	}
}
