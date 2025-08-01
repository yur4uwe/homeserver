package handlers

import (
	"fmt"
	"homeserver/pkg/globals"
	"homeserver/pkg/schemas"
	"homeserver/pkg/templates"
	"math"
	"net/url"
	"strings"
)

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

func FileToSearchEntry(file *schemas.UserFile) templates.SearchEntry {
	entry := templates.SearchEntry{
		ID:        file.ID,
		Name:      file.Filename,
		Path:      strings.Join(strings.Split(file.Path, "/")[1:], "/"),
		Size:      ParseSize(file.Size),
		Extension: "." + file.Extension,
		IsFolder:  false,
	}

	if entry.Path != "" {
		entry.URIPath = ParseQuerySlice(strings.Split(entry.Path, "/"))
		entry.Path += "/"
		entry.URIPath += "/"
	}

	entry.URIName = url.QueryEscape(entry.Name)

	extension := file.Extension
	if globals.IsFileTypeAvailable(extension) {
		entry.Type = "icons/" + extension
	} else {
		entry.Type = "file"
	}

	return entry
}

func FolderToSearchEntry(folder *schemas.UserFolder) templates.SearchEntry {
	entry := templates.SearchEntry{
		ID:       folder.ID,
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
	return entry
}

func TransformEntries(file_entries []*schemas.UserFile, folder_entries []*schemas.UserFolder) []templates.SearchEntry {
	var entries []templates.SearchEntry

	for _, file := range file_entries {
		entry := FileToSearchEntry(file)
		entries = append(entries, entry)
	}

	for _, folder := range folder_entries {
		entry := FolderToSearchEntry(folder)
		entries = append(entries, entry)
	}

	return entries
}
