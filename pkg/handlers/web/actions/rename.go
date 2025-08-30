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
	"net/url"
	"os"
)

var (
	NewAlert      = parser.NewAlert
	ErrorResponse = errresp.ErrorResponse
)

func RenameFileHandler(w http.ResponseWriter, r *http.Request, repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext, fileId int) {
	is_folder := r.URL.Query().Get("is_folder") == "true"

	orig_name, err := url.QueryUnescape(r.URL.Query().Get("orig_name"))
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("Error unescaping file path: %v", err), http.StatusBadRequest)
		return
	}

	desiredName := r.URL.Query().Get("desired_name")
	if desiredName == "" {
		ErrorResponse(w, "Desired name parameter is required for rename action", http.StatusBadRequest)
		return
	}

	if is_folder {

	} else {
		oldFile := fschema.New().
			WithID(fileId).
			WithFilename(orig_name).
			WithOwnerID(reqCont.Requester.ID)

		fvp := fschema.New().WithFilename(desiredName)

		f, err := repos.FileRepo.UpdateFileEntry(r.Context(), oldFile, fvp, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, fmt.Sprintf("Error renaming file: %v", err), http.StatusInternalServerError)
			return
		}

		oldFile.WithExtension(f.Extension).WithPath(f.Path)

		if err := os.Rename(oldFile.AbsolutePath(), f.AbsolutePath()); err != nil {
			ErrorResponse(w, "Error renaming file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", formatters.AlertEventBuilder("Success!", "File successfully renamed", constants.Success))
		w.Write([]byte(f.Filename + "." + f.Extension))
	}
}
