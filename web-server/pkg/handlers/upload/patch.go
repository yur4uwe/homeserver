package upload

import (
	"homeserver/constants/fconsts"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/web-server/pkg/middleware"
	"net/http"
)

func PatchFile(repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Patch file request received")

		newFileName := r.URL.Query().Get("new_name")
		oldFileName := r.URL.Query().Get("old_name")

		if newFileName == "" || oldFileName == "" {
			reqCont.LogHistory.ErrLog("Missing filename parameter", "PatchFile")
			ErrorResponse(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		file := fschema.NewUserFileBuilder().
			WithFullName(oldFileName).
			WithPath(reqCont.Requester.Username).
			WithOwnerID(reqCont.Requester.ID)

		err := repos.FileRepo.UpdateFileEntry(r.Context(), file, []string{fconsts.Filename}, reqCont.Requester)
		if err != nil {
			reqCont.LogHistory.ErrLog("Error renaming file: "+err.Error(), "PatchFile")
			ErrorResponse(w, "Error renaming file", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("File Renamed successfully: " + oldFileName + " to " + newFileName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(file.Filename + "." + file.Extension))
	}
}
