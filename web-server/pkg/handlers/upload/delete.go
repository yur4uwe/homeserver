package upload

import (
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/utils"
	"net/http"
)

var ErrorResponse = utils.ErrorResponse

func DeleteFile(repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Delete file request received")

		filePath := r.URL.Query().Get("filename")
		if filePath == "" {
			reqCont.LogHistory.ErrLog("Missing filename parameter", "DeleteFile")
			ErrorResponse(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		file := fschema.NewUserFileBuilder().
			WithPath(filePath).
			WithOwnerID(reqCont.Requester.ID).
			WithFullName(filePath)

		if _, err := repos.FileRepo.DeleteFile(r.Context(), file, reqCont.Requester); err != nil {
			reqCont.LogHistory.ErrLog("Error deleting file: "+err.Error(), "DeleteFile")
			ErrorResponse(w, "Error deleting file", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("File successfully deleted: " + filePath)
		w.WriteHeader(http.StatusOK)
	}
}
