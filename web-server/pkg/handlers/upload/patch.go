package upload

import (
	"homeserver/db/repositories"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/file_schema"
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

		oldFile := file_schema.NewUserFileBuilder().
			WithFullName(oldFileName).
			WithPath(reqCont.Requester.Username).
			WithOwnerID(reqCont.Requester.ID)

		newFilePath, err := connection.RenameFile(r.Context(), oldFile, newFileName, reqCont.Requester)
		if err != nil {
			reqCont.LogHistory.ErrLog("Error renaming file: "+err.Error(), "PatchFile")
			ErrorResponse(w, "Error renaming file", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("File Renamed successfully: " + oldFileName + " to " + newFileName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(newFilePath))
	}
}
