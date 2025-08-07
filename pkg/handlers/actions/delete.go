package actions

import (
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func DeleteFileHandler(w http.ResponseWriter, r *http.Request, connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string, is_folder bool) {
	filePath := strings.Join(append(relativePath, fileName), "/")

	file := file_schema.NewUserFileBuilder().
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
