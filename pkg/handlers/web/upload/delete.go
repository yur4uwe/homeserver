package upload

import (
	"homeserver/db"
	middleware "homeserver/pkg/middleware/web"
	"net/http"
	"strings"
)

type UploadFileHandler struct{}

func DeleteFile(repos *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Delete file request received")

		filePath := r.URL.Query().Get("filename")
		if filePath == "" {
			reqCont.LogHistory.ErrLog("Missing filename parameter", "DeleteFile")
			ErrorResponse(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		parts := strings.Split(filePath, ".")
		filename := strings.Join(parts[:len(parts)-1], ".")
		extension := parts[len(parts)-1]

		file := &db.UserFile{
			Path:      filePath, // Keeping original behavior where Path = filePath
			Filename:  filename,
			Extension: extension,
			OwnerID:   reqCont.Requester.ID,
		}

		if _, err := repos.Files.Delete(r.Context(), file, reqCont.Requester); err != nil {
			reqCont.LogHistory.ErrLog("Error deleting file: "+err.Error(), "DeleteFile")
			ErrorResponse(w, "Error deleting file", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("File successfully deleted: " + filePath)
		w.WriteHeader(http.StatusOK)
	}
}
