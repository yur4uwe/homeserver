package handlers

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
)

func IsInActiveDirectory(directory []string, path []string, is_folder bool) bool {
	if !is_folder && len(path) != len(directory) {
		return false
	}

	if is_folder && len(path) != len(directory)+1 {
		return false
	}

	for i := 0; i < len(directory); i++ {
		if directory[i] != path[i] {
			return false
		}
	}
	return true
}

func ParseQuerySlice(query []string) string {
	new_query := make([]string, len(query))

	for i, v := range query {
		new_query[i] = url.QueryEscape(v)
	}
	return strings.Join(new_query, "/")
}

func ParseSize(size int64) string {
	size_float := float64(size)
	size_units := []string{"B", "KB", "MB", "GB", "TB"}

	i := 0
	for size_float >= math.Pow(1024, float64(i+1)) && i < len(size_units)-1 {
		i++
	}

	return fmt.Sprintf("%.2f %s", size_float/math.Pow(1024, float64(i)), size_units[i])
}

func HasFolder(folder_name string, entries []templates.SearchEntry) bool {
	for _, entry := range entries {
		if entry.IsFolder && entry.Name == folder_name {
			return true
		}
	}
	return false
}

func TransformEntries(file_entries []*db.UserFile, folder_entries []*db.UserFolder) []templates.SearchEntry {
	var entries []templates.SearchEntry

	for _, file := range file_entries {
		entry := templates.SearchEntry{
			Name:     file.Filename,
			Path:     strings.Join(strings.Split(file.Path, "/")[1:], "/"),
			Size:     ParseSize(file.Size),
			IsFolder: false,
		}

		entry.URIName = url.QueryEscape(entry.Name)
		entry.URIPath = ParseQuerySlice(strings.Split(entry.Path, "/"))

		extension := strings.Split(file.Filename, ".")[1]
		if globals.IsFileTypeAvailable(extension) {
			entry.Type = "icons/" + extension
		} else {
			entry.Type = "file"
		}

		entries = append(entries, entry)
	}

	for _, folder := range folder_entries {
		entry := templates.SearchEntry{
			Name:     folder.Name,
			Path:     strings.Join(strings.Split(folder.Path, "/")[1:], "/"),
			Size:     "",
			Type:     "folder",
			IsFolder: true,
		}

		if entry.Path != "" {
			entry.URIPath = ParseQuerySlice(strings.Split(entry.Path, "/"))
			entry.Path += "/"
			entry.URIPath += "/"
		}
		entry.URIName = url.QueryEscape(entry.Name)

		entries = append(entries, entry)
	}

	return entries
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
		directory_string, err := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/search/"))
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API url.QueryUnescape error: "+err.Error(), "SearchHandler:30")
			http.Error(w, "Internal Server Error: Failed to unescape directory", http.StatusInternalServerError)
			return
		}

		directory := []string{reqCont.Requester.Username}
		if directory_string != "" {
			directory = append(directory, strings.Split(directory_string, "/")...)
		}

		remote_path := strings.TrimPrefix(r.Header.Get("Referer"), fmt.Sprintf("http://%s/search/", r.Host))

		directory_without_username := directory[1:]

		if remote_path != strings.Join(directory_without_username, "/") {
			reqCont.LogHistory.Log(fmt.Sprintf("[Handler] Restoring directory for client: '%s'", directory))
			if len(directory) > 0 {
				directory_string = ParseQuerySlice(directory_without_username)
			}

			w.Header().Set("HX-Push-Url", "/search/"+directory_string)
		}

		reqCont.LogHistory.Log(fmt.Sprintf("[Handler] Search API query: '%s', directory: '%s'", query, strings.Join(directory, "/")))

		file_entries, err := connection.SearchFiles(r.Context(), query, reqCont.Requester, directory)
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			http.Error(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
			return
		}

		folder_entries, err := connection.ExtractFoldersFromDirectory(r.Context(), directory, reqCont.Requester)
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			http.Error(w, "Internal Server Error: Failed to search folders", http.StatusInternalServerError)
			return
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

		var parent []string
		if len(directory_without_username) > 0 {
			parent = directory_without_username[:len(directory_without_username)-1]
		}

		data := &templates.SearchPageData{
			Query:        query,
			Parent:       strings.Join(parent, "/"),
			Directory:    strings.Join(directory_without_username, "/"),
			URIQuery:     url.QueryEscape(query),
			URIParent:    ParseQuerySlice(parent),
			URIDirectory: ParseQuerySlice(directory_without_username),
			Entries:      relevant_entries,
		}

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
