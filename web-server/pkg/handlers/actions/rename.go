package actions

import (
	"fmt"
	"homeserver/constants"
	"homeserver/constants/fconsts"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/parser"
	"homeserver/web-server/pkg/utils/err"
	"homeserver/web-server/pkg/utils/formatters"
	"net/http"
	"strings"
)

var (
	NewAlert      = parser.NewAlert
	ErrorResponse = err.ErrorResponse
)

func RenameFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, relativePath []string, fileName, desiredName string, is_folder bool) {
	if desiredName == "" || fileName == "" {
		ErrorResponse(w, "Desired name parameter is required for rename action", http.StatusBadRequest)
		return
	}

	oldFile := fschema.NewUserFileBuilder().
		WithFullName(fileName).
		WithPath(strings.Join(relativePath, "/")).
		WithOwnerID(reqCont.Requester.ID)

	fvp := map[string]any{
		fconsts.Filename: desiredName,
	}

	err := repos.FileRepo.UpdateFileEntry(r.Context(), oldFile, fvp, reqCont.Requester)
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("Error renaming file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", formatters.AlertEventBuilder("Success!", "File successfully renamed", constants.Success))
	w.Write([]byte(oldFile.Filename + "." + oldFile.Extension))
}
