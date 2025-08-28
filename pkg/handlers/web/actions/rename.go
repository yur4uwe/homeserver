package actions

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils/errresp"
	"homeserver/pkg/utils/formatters"
	"net/http"
	"strings"
)

var (
	NewAlert      = parser.NewAlert
	ErrorResponse = errresp.ErrorResponse
)

func RenameFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName string, is_folder bool) {
	if fileName == "" {
		ErrorResponse(w, "File name parameter is required for rename action", http.StatusBadRequest)
		return
	}

	desiredName := r.URL.Query().Get("desired_name")
	if desiredName == "" {
		ErrorResponse(w, "Desired name parameter is required for rename action", http.StatusBadRequest)
		return
	}

	oldFile := fschema.NewUserFileBuilder().
		WithFullName(fileName).
		WithPath(strings.Join(relativePath, "/")).
		WithOwnerID(reqCont.Requester.ID)

	fvp := fschema.NewUserFileBuilder().WithFilename(desiredName)

	f, err := repos.FileRepo.UpdateFileEntry(r.Context(), oldFile, fvp, reqCont.Requester)
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("Error renaming file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", formatters.AlertEventBuilder("Success!", "File successfully renamed", constants.Success))
	w.Write([]byte(f.Filename + "." + f.Extension))
}
