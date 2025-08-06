package upload

import (
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas/file_schema"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func DownloadFile(connection *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Download file request received")

		// Monitor for cancellation
		ctx := r.Context()
		done := make(chan struct{})
		go func() {
			<-ctx.Done()
			reqCont.LogHistory.ErrLog("Request canceled by client", "DownloadFile")
			done <- struct{}{}
		}()

		var file_sig []string

		multipart, header, err := r.FormFile("file")
		if err != nil {
			select {
			case <-done:
				// Handle cancellation
				reqCont.LogHistory.ErrLog("Request canceled during file upload", "DownloadFile")
				utils.ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
				return
			default:
				// Handle other errors
				reqCont.LogHistory.ErrLog("Error retrieving file from form: "+err.Error(), "DownloadFile")
				utils.ErrorResponse(w, "Error retrieving file", http.StatusBadRequest)
				return
			}
		}
		defer multipart.Close()

		path_str := r.FormValue("folder")
		if path_str != "" {
			path_segments := strings.Split(path_str, "/")
			path_str = reqCont.Requester.Username + "/" + strings.Join(path_segments[:len(path_segments)-1], "/")
			file_sig = strings.Split(path_segments[len(path_segments)-1], ".")
		} else {
			file_sig = strings.Split(header.Filename, ".")
			path_str = reqCont.Requester.Username
		}

		new_file := file_schema.NewUserFileBuilder().
			WithFilename(file_sig[0]).
			WithPath(path_str).
			WithExtension(file_sig[1]).
			WithSize(header.Size).
			WithContentType(header.Header.Get("Content-Type")).
			WithOwnerID(reqCont.Requester.ID)

		_, saved_name, err := connection.SaveFile(r.Context(), new_file, &multipart, reqCont.Requester)
		if err != nil {
			select {
			case <-done:
				// Handle cancellation
				reqCont.LogHistory.ErrLog("Request canceled during file saving", "DownloadFile")
				utils.ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
				return
			default:
				// Handle other errors
				reqCont.LogHistory.ErrLog("Error saving file: "+err.Error(), "DownloadFile")
				utils.ErrorResponse(w, "Error saving file", http.StatusInternalServerError)
				return
			}
		}

		reqCont.LogHistory.Log("File successfully processed for upload")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(saved_name))
	}
}
