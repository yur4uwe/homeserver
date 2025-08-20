package handlers

import (
	"fmt"
	"homeserver/constants"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"net/http"
)

func PageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Welcome to the Home Server</h1>"))
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	entry := r.Context().Value(constants.EntryKey).(*fschema.UserFile)
	if entry == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	fmt.Printf(
		"Requested file: %s\nExtension: %s\nPath: %s\nOwner: %d\n",
		entry.Filename, entry.Extension, entry.Path, entry.OwnerID,
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>File Management</h1>"))
}

func FolderHandler(w http.ResponseWriter, r *http.Request) {
	entry := r.Context().Value(constants.EntryKey).(*folschema.UserFolder)
	if entry == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	fmt.Printf(
		"Requested folder: %s\nPath: %s\nOwner: %d\n",
		entry.Name, entry.Path, entry.OwnerID,
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Folder Management</h1>"))
}
