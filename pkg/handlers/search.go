package handlers

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
)

func IsFileInActiveDirectory(directory []string, path []string) bool {
	if len(path) != len(directory) {
		return false
	}

	for i := 0; i < len(directory); i++ {
		if directory[i] != path[i] {
			return false
		}
	}
	return true
}

func IsFolderInActiveDirectory(directory []string, path []string) bool {
	if len(path) != len(directory)+1 {
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

func FilterRelevantEntries(files []*db.UserFiles, directory []string, base []string) []templates.SearchEntry {
	// directory = strings.TrimPrefix(directory, "/")
	fmt.Println("Directory:", strings.Join(directory, "/"))
	fmt.Println("Directory Length:", len(directory))
	fmt.Println("Base:", strings.Join(base, "/"))

	var filteredEntries []templates.SearchEntry
	for _, file := range files {
		normalizedPath := strings.TrimPrefix(strings.ReplaceAll(file.Path, "\\", "/"), strings.Join(base, "/")+"/") // Normalize path separators

		relative_path := strings.Split(normalizedPath, "/")

		relative_path = relative_path[:len(relative_path)-1] // Remove the file name

		fmt.Println("Relative Path:", strings.Join(relative_path, "/"))
		fmt.Println("is file in active directory:", IsFileInActiveDirectory(directory, relative_path))
		fmt.Println("is folder in active directory:", IsFolderInActiveDirectory(directory, relative_path))

		// is a folder on this level
		if len(relative_path)-len(directory) > 0 {
			folder := templates.SearchEntry{
				IsFolder: true,
				Name:     relative_path[len(relative_path)-1],
				Path:     strings.Join(relative_path[:len(relative_path)-1], "/"),
				Type:     "folder",
				Size:     "",
			}

			folder.URIName = url.QueryEscape(folder.Name)
			folder.URIPath = ParseQuerySlice(relative_path[:len(relative_path)-1])

			if folder.Path != "" {
				folder.Path += "/"
				folder.URIPath += "/"
			}

			if !IsFolderInActiveDirectory(directory, relative_path) {
				continue
			}

			if !HasFolder(folder.Name, filteredEntries) {
				filteredEntries = append(filteredEntries, folder)
			}
			continue
		}

		entry := templates.SearchEntry{
			Name:     file.Filename,
			Path:     strings.Join(relative_path, "/"),
			Size:     ParseSize(file.Size),
			IsFolder: false,
		}

		entry.URIName = url.QueryEscape(entry.Name)
		entry.URIPath = ParseQuerySlice(relative_path)

		extension := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]

		if globals.IsFileTypeAvailable(extension) {
			entry.Type = "icons/" + extension
		} else {
			entry.Type = "file"
		}

		if !IsFileInActiveDirectory(directory, relative_path) {
			continue
		}
		filteredEntries = append(filteredEntries, entry)
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
		directory_string, err := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/search/"))
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API url.QueryUnescape error: "+err.Error(), "SearchHandler:30")
			http.Error(w, "Internal Server Error: Failed to unescape directory", http.StatusInternalServerError)
			return
		}

		var directory []string
		if directory_string == "" {
			directory = []string{}
		} else {
			directory = strings.Split(directory_string, "/")
		}

		if strings.TrimPrefix(r.Header.Get("Referer"), fmt.Sprintf("http://%s/search/", r.Host)) != strings.Join(directory, "/") {
			reqCont.LogHistory.Log(fmt.Sprintf("[Handler] Restoring directory for client: '%s'", directory))
			if len(directory) > 0 {
				directory_string = ParseQuerySlice(directory)
			}

			fmt.Println("Directory:", directory)
			fmt.Println("Directory Length:", len(directory))
			fmt.Println("Directory String:", directory_string)

			w.Header().Set("HX-Push-Url", "/search/"+directory_string)
		}

		reqCont.LogHistory.Log(fmt.Sprintf("[Handler] Search API query: '%s', directory: '%s'", query, strings.Join(directory, "/")))

		entries, err := connection.SearchFiles(r.Context(), query, reqCont.AuthToken)
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			http.Error(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
			return
		}

		relevant_entries := FilterRelevantEntries(
			entries,
			directory,
			strings.Split(os.Getenv(globals.EnvStoragePath)+reqCont.Requester.Username, "/"),
		)

		sort.Slice(relevant_entries, func(i, j int) bool {
			if relevant_entries[i].IsFolder && !relevant_entries[j].IsFolder {
				return true
			} else if !relevant_entries[i].IsFolder && relevant_entries[j].IsFolder {
				return false
			}
			return relevant_entries[i].Name < relevant_entries[j].Name
		})

		var parent []string
		if len(directory) > 0 {
			parent = directory[:len(directory)-1]
		}

		data := &templates.SearchPageData{
			Query:        query,
			Parent:       strings.Join(parent, "/"),
			Directory:    strings.Join(directory, "/"),
			URIQuery:     url.QueryEscape(query),
			URIParent:    ParseQuerySlice(parent),
			URIDirectory: ParseQuerySlice(directory),
			Entries:      relevant_entries,
		}

		reqCont.LogHistory.Log(fmt.Sprintln("+-----+\nData:", data, "\n+-----+"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		var tmpl *template.Template
		if reqCont.ToRestorePage || reqCont.RerenderBody {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_base.html", "../static/templates/search_results.html"))

			tmpl.ExecuteTemplate(w, "search_base.html", data)
		} else {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_results.html"))

			tmpl.Execute(w, data)
		}
	}
}
