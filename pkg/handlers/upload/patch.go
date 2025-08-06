package upload

import (
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/utils"
	"net/http"
)

func PatchFile(connection *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Patch file request received")

		newFileName := r.URL.Query().Get("new_name")
		oldFileName := r.URL.Query().Get("old_name")

		newFilePath, err := connection.RenameFile(r.Context(), oldFileName, newFileName, reqCont.Requester)
		if err != nil {
			reqCont.LogHistory.ErrLog("Error renaming file: "+err.Error(), "PatchFile")
			utils.ErrorResponse(w, "Error renaming file", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("File Renamed successfully: " + oldFileName + " to " + newFileName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(newFilePath))
	}
}
