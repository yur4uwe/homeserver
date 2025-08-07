package upload

import (
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/utils"
	"net/http"
	"path/filepath"
	"strings"
)

func DeleteFile(connection interfaces.DatabaseInterface, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Delete file request received")

		filePath := r.URL.Query().Get("filename")
		if filePath == "" {
			reqCont.LogHistory.ErrLog("Missing filename parameter", "DeleteFile")
			utils.ErrorResponse(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		file := file_schema.NewUserFileBuilder().
			WithPath(filePath).
			WithOwnerID(reqCont.Requester.ID).
			WithFilename(strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))).
			WithExtension(filepath.Ext(filePath)[1:])

		if _, err := connection.DeleteFile(r.Context(), file, reqCont.Requester.Username); err != nil {
			reqCont.LogHistory.ErrLog("Error deleting file: "+err.Error(), "DeleteFile")
			utils.ErrorResponse(w, "Error deleting file", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("File successfully deleted: " + filePath)
		w.WriteHeader(http.StatusOK)
	}
}
