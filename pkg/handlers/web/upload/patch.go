package upload

import (
	"homeserver/db"
	middleware "homeserver/pkg/middleware/web"
	"net/http"
	"net/url"
	"strings"
)

func PatchFile(repos *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Patch file request received")

		newFileName, err := url.QueryUnescape(r.URL.Query().Get("new_name"))
		if err != nil {
			reqCont.LogHistory.ErrLog("Error unescaping new filename: "+err.Error(), "PatchFile")
			ErrorResponse(w, "Error unescaping new filename", http.StatusBadRequest)
			return
		}
		oldFileName, err := url.QueryUnescape(r.URL.Query().Get("old_name"))
		if err != nil {
			reqCont.LogHistory.ErrLog("Error unescaping old filename: "+err.Error(), "PatchFile")
			ErrorResponse(w, "Error unescaping old filename", http.StatusBadRequest)
			return
		}

		if newFileName == "" || oldFileName == "" {
			reqCont.LogHistory.ErrLog("Missing filename parameter", "PatchFile")
			ErrorResponse(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		parts := strings.Split(oldFileName, ".")
		oldFilename := strings.Join(parts[:len(parts)-1], ".")
		oldExtension := parts[len(parts)-1]

		file := &db.UserFile{
			Filename:  oldFilename,
			Extension: oldExtension,
			Path:      reqCont.Requester.Username,
			OwnerID:   reqCont.Requester.ID,
		}

		sig := strings.Split(newFileName, ".")

		updates := &db.UserFile{
			Filename: sig[0],
		}

		f, err := repos.Files.Update(r.Context(), file, updates, reqCont.Requester)
		if err != nil {
			reqCont.LogHistory.ErrLog("Error renaming file: "+err.Error(), "PatchFile")
			ErrorResponse(w, "Error renaming file", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("File Renamed successfully: " + oldFileName + " to " + newFileName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(f.Filename + "." + f.Extension))
	}
}
