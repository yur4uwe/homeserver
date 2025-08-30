package actions

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	middleware "homeserver/pkg/middleware/web"
	"net/http"
	"os"
)

func DeleteFileHandler(
	w http.ResponseWriter, r *http.Request,
	repos *repositories.Repositories,
	reqCont *middleware.AdditionalRequestContext,
	entryId int,
) {
	isFolder := r.URL.Query().Get("is_folder") == "true"

	var alert string
	if isFolder {
		folder := folschema.New().WithID(entryId).WithOwnerID(reqCont.Requester.ID)

		fol, err := repos.FolderRepo.DeleteFolder(r.Context(), folder, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, fmt.Sprintf("Error deleting folder: %v", err), http.StatusInternalServerError)
			return
		}

		alert = NewAlert(
			"Folder Deleted",
			"check",
			"Folder has been deleted successfully.",
			constants.Success,
		)

		if err := os.RemoveAll(fol.AbsolutePath()); err != nil {
			ErrorResponse(w, fmt.Sprintf("Error removing folder from filesystem: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", fmt.Sprintf("delete-success-folder-%d", fol.ID))
	} else {
		file := fschema.New().WithID(entryId).WithOwnerID(reqCont.Requester.ID)

		f, err := repos.FileRepo.DeleteFile(r.Context(), file, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, fmt.Sprintf("Error deleting file: %v", err), http.StatusInternalServerError)
			return
		}

		alert = NewAlert(
			"File Deleted",
			"check",
			"File has been deleted successfully.",
			constants.Success,
		)

		if err := os.RemoveAll(f.AbsolutePath()); err != nil {
			ErrorResponse(w, fmt.Sprintf("Error removing file from filesystem: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", fmt.Sprintf("delete-success-file-%d", f.ID))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(alert))
}
