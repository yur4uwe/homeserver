package actions

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	folschema "homeserver/db/repositories/folder_repo/schema"
	"homeserver/web-server/pkg/middleware"
	"net/http"
	"strings"
)

func DeleteFileHandler(
	w http.ResponseWriter, r *http.Request,
	repos *repositories.Repositories,
	reqCont *middleware.AdditionalRequestContext,
	relativePath []string, fileName string,
) {
	entryNameSlice := strings.Split(fileName, ".")

	var alert string
	if len(entryNameSlice) < 2 {
		folder := folschema.NewUserFolderBuilder().
			WithName(entryNameSlice[0]).
			WithPath(strings.Join(relativePath, "/")).
			WithOwnerID(reqCont.Requester.ID)

		id, err := repos.FolderRepo.DeleteFolder(r.Context(), folder, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, fmt.Sprintf("Error deleting folder: %v", err), http.StatusInternalServerError)
			return
		}

		alert = NewAlert(
			"Folder Deleted",
			"check",
			fmt.Sprintf("Folder '%s' has been deleted successfully.", folder.Name),
			constants.Success,
		)

		w.Header().Set("HX-Trigger", fmt.Sprintf("delete-success-folder-%d", id))
	} else {
		file := fschema.NewUserFileBuilder().
			WithFilename(entryNameSlice[0]).
			WithExtension(entryNameSlice[1]).
			WithPath(strings.Join(relativePath, "/")).
			WithOwnerID(reqCont.Requester.ID)

		id, err := repos.FileRepo.DeleteFile(r.Context(), file, reqCont.Requester)
		if err != nil {
			ErrorResponse(w, fmt.Sprintf("Error deleting file: %v", err), http.StatusInternalServerError)
			return
		}

		alert = NewAlert(
			"File Deleted",
			"check",
			fmt.Sprintf("File '%s.%s' has been deleted successfully.", entryNameSlice[0], entryNameSlice[1]),
			constants.Success,
		)

		w.Header().Set("HX-Trigger", fmt.Sprintf("delete-success-file-%d", id))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(alert))
}
