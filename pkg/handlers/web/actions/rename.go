package actions

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils"
	"net/http"
	"net/url"
	"os"
)

var (
	NewAlert      = parser.NewAlert
	ErrorResponse = utils.ErrorResponse
)

func RenameFileHandler(w http.ResponseWriter, r *http.Request, repos *db.Database, reqCont *middleware.AdditionalRequestContext, fileId int) {
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
		oldFile := &db.UserFile{
			ID:       fileId,
			Filename: orig_name,
			OwnerID:  reqCont.Requester.ID,
		}

		updates := &db.UserFile{Filename: desiredName}

		f, err := repos.Files.Update(r.Context(), oldFile, updates, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, fmt.Sprintf("Error renaming file: %v", err), http.StatusInternalServerError)
			return
		}

		oldFile.Extension = f.Extension
		oldFile.Path = f.Path

		if err := os.Rename(oldFile.AbsolutePath(), f.AbsolutePath()); err != nil {
			ErrorResponse(w, "Error renaming file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", utils.AlertEventBuilder("Success!", "File successfully renamed", constants.Success))
		w.Write([]byte(f.Filename + "." + f.Extension))
	}
}
