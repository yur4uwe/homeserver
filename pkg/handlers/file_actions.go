package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/schemas"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func DeleteFileHandler(w http.ResponseWriter, r *http.Request, connection *db.Database, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string, is_folder bool) {
	filePath := strings.Join(append(relativePath, fileName), "/")

	file := schemas.NewUserFileBuilder().
		WithFilename(fileName).
		WithPath(filePath).
		WithOwnerID(reqCont.Requester.ID)

	id, err := connection.DeleteFile(r.Context(), file, reqCont.Requester.Username)
	if err != nil {
		utils.ErrorResponse(w, fmt.Sprintf("Error deleting file: %v", err), http.StatusInternalServerError)
		return
	}

	alert := parser.NewAlert(
		"File Deleted",
		"check",
		fmt.Sprintf("File '%s' has been deleted successfully.", fileName),
		utils.AlertSuccess,
	)

	entry_type := "file"
	if is_folder {
		entry_type = "folder"
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("delete-success-%s-%d", entry_type, id))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Retarget", "#modal-content")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Write([]byte(alert))
}

func RenameFileHandler(w http.ResponseWriter, r *http.Request, connection *db.Database, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName, desiredName string, is_folder bool) {
	if desiredName == "" {
		utils.ErrorResponse(w, "Desired name parameter is required for rename action", http.StatusBadRequest)
		return
	}

	filePath := strings.Join(append(relativePath, fileName), "/")

	newName, err := connection.RenameFile(r.Context(), filePath, desiredName, reqCont.Requester)
	if err != nil {
		utils.ErrorResponse(w, fmt.Sprintf("Error renaming file: %v", err), http.StatusInternalServerError)
		return
	}

	alert := parser.NewAlert(
		"File Renamed",
		"check",
		fmt.Sprintf("File '%s' has been renamed to '%s'.", fileName, newName),
		utils.AlertSuccess,
	)

	w.Header().Set("HX-Trigger", fmt.Sprintf("rename-success-%s-%s", strings.Join(relativePath, "/"), fileName))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Retarget", "#modal-content")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Write([]byte(alert))
}

func OpenFileHandler(w http.ResponseWriter, r *http.Request, connection *db.Database, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string) {
	filePath := strings.Join(append(relativePath, fileName), "/")

	file, err := connection.GetFile(r.Context(), filePath, reqCont.Requester)
	if err != nil {
		utils.ErrorResponse(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, file.Path)
}
