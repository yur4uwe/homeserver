package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"html/template"
	"net/http"
	"sort"
	"strings"
)

func HasFolder(folder_name string, entries []templates.SearchEntry) bool {
	for _, entry := range entries {
		if entry.IsFolder && entry.Name == folder_name {
			return true
		}
	}
	return false
}

func FilterRelevantEntries(entries []templates.SearchEntry, directory string) []templates.SearchEntry {
	directory = strings.TrimSuffix(directory, "/")
	// directory = strings.TrimPrefix(directory, "/")

	var filteredEntries []templates.SearchEntry
	for _, entry := range entries {

		relative_path := strings.TrimPrefix(entry.Path, directory+"/")

		path_slice := strings.Split(relative_path, "/")
		if len(path_slice) == 2 {
			folder_entry := templates.SearchEntry{
				Name:     path_slice[0],
				Path:     directory,
				IsFolder: true,
				Size:     "",
			}

			if !HasFolder(folder_entry.Name, filteredEntries) {
				filteredEntries = append(filteredEntries, folder_entry)
			}
		} else if len(path_slice) == 1 {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	return filteredEntries
}

func SearchHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Search API accessed")

		query := r.URL.Query().Get("query")
		directory := strings.TrimPrefix(r.URL.Path, "/search/")
		if directory != "" && strings.TrimPrefix(r.Header.Get("Referer"), fmt.Sprintf("http://%s/search/", r.Host)) != directory {
			w.Header().Set("HX-Push-Url", "/search/"+directory)
		}

		reqCont.LogHistory.Log(fmt.Sprintf("[Handler] Search API query: '%s', directory: '%s'", query, directory))

		entries, err := connection.SearchFiles(r.Context(), query, reqCont.AuthToken)
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			http.Error(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
			return
		}

		// relevant_entries := FilterRelevantEntries(entries, os.Getenv(globals.EnvStoragePath)+directory)
		relevant_entries := FilterRelevantEntries(entries, "/"+directory)

		sort.Slice(relevant_entries, func(i, j int) bool {
			if relevant_entries[i].IsFolder && !relevant_entries[j].IsFolder {
				return true
			} else if !relevant_entries[i].IsFolder && relevant_entries[j].IsFolder {
				return false
			}
			return relevant_entries[i].Name < relevant_entries[j].Name
		})

		parent := strings.Join(strings.Split(directory, "/")[:len(strings.Split(directory, "/"))-1], "/")

		data := &templates.SearchPageData{
			Query:     query,
			Parent:    parent,
			Directory: directory,
			Entries:   relevant_entries,
		}

		reqCont.LogHistory.Log(fmt.Sprintln("+-----+\nData:", data, "\n+-----+"))

		tmpl := template.Must(template.ParseFiles("../static/templates/search.html"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}
