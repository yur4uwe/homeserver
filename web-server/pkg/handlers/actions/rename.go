package actions

import (
	"fmt"
	"homeserver/db/repositories"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

var (
	NewAlert      = parser.NewAlert
	ErrorResponse = utils.ErrorResponse
)

func RenameFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName, desiredName string, is_folder bool) {
	if desiredName == "" || fileName == "" {
		ErrorResponse(w, "Desired name parameter is required for rename action", http.StatusBadRequest)
		return
	}

	oldFile := file_schema.NewUserFileBuilder().
		WithFullName(fileName).
		WithPath(strings.Join(relativePath, "/")).
		WithOwnerID(reqCont.Requester.ID)

	newName, err := connection.RenameFile(r.Context(), oldFile, desiredName, reqCont.Requester)
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("Error renaming file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", `{"raise-alert": {"title": "Success!", "value": "File successfully renamed", "status": "success", "icon": "check"}}`)
	w.Write([]byte(newName))
}
