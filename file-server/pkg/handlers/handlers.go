package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"homeserver/constants"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/file-server/internal/respfail"
	"homeserver/file-server/internal/zipper"
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
		respfail.ErrorResponse(w, r, "Not Found", http.StatusNotFound)
		return
	}

	fmt.Printf(
		"Requested file: %s\nExtension: %s\nPath: %s\nOwner: %d\n",
		entry.Filename, entry.Extension, entry.Path, entry.OwnerID,
	)

	r = r.WithContext(context.WithValue(r.Context(), "status", http.StatusOK))

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", entry.Filename))
	w.Header().Set("Content-Type", "application/zip")
	w.WriteHeader(http.StatusOK)

	zw := zip.NewWriter(w)
	defer zw.Close()

	if err := zipper.ZipFile(zw, entry.AbsolutePath(), entry.Filename+"."+entry.Extension); err != nil {
		respfail.ErrorResponse(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
}

func FolderHandler(w http.ResponseWriter, r *http.Request) {
	entry := r.Context().Value(constants.EntryKey).(*folschema.UserFolder)
	if entry == nil {
		respfail.ErrorResponse(w, r, "Not Found", http.StatusNotFound)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "status", http.StatusOK))

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", entry.Name))
	w.Header().Set("Content-Type", "application/zip")
	w.WriteHeader(http.StatusOK)

	if err := zipper.ZipFolder(w, entry); err != nil {
		respfail.ErrorResponse(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
}
