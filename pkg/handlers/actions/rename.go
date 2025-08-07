package actions

import (
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func RenameFileHandler(w http.ResponseWriter, r *http.Request, connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName, desiredName string, is_folder bool) {
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
