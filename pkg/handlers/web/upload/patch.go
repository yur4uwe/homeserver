package upload

import (
	"homeserver/constants/fconsts"
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	middleware "homeserver/pkg/middleware/web"
	"net/http"
	"net/url"
	"strings"
)

func PatchFile(repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
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

		file := fschema.NewUserFileBuilder().
			WithFullName(oldFileName).
			WithPath(reqCont.Requester.Username).
			WithOwnerID(reqCont.Requester.ID)

		sig := strings.Split(newFileName, ".")

		fvp := map[string]any{
			fconsts.Filename: sig[0],
		}

		err = repos.FileRepo.UpdateFileEntry(r.Context(), file, fvp, reqCont.Requester)
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
