package upload

import (
	"homeserver/db/repositories"
	fschema "homeserver/db/repositories/file_repo/schema"
	"homeserver/web-server/pkg/middleware"
	"net/http"
	"strings"
)

func DownloadFile(repos *repositories.Repositories, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
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
				reqCont.LogHistory.ErrLog("Request canceled during file upload", "DownloadFile:31")
				ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
				return
			default:
				// Handle other errors
				reqCont.LogHistory.ErrLog("Error retrieving file from form: "+err.Error(), "DownloadFile:36")
				ErrorResponse(w, "Error retrieving file", http.StatusBadRequest)
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

		new_file := fschema.NewUserFileBuilder().
			WithFilename(file_sig[0]).
			WithPath(path_str).
			WithExtension(file_sig[1]).
			WithSize(header.Size).
			WithContentType(header.Header.Get("Content-Type")).
			WithOwnerID(reqCont.Requester.ID)

		inserted_file, err := repos.FileRepo.CreateFile(r.Context(), new_file, reqCont.Requester)
		if err != nil {
			select {
			case <-done:
				// Handle cancellation
				reqCont.LogHistory.ErrLog("Request canceled during file saving", "DownloadFile:66")
				ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
				return
			default:
				// Handle other errors
				reqCont.LogHistory.ErrLog("Error saving file: "+err.Error(), "DownloadFile:71")
				ErrorResponse(w, "Error saving file", http.StatusInternalServerError)
				return
			}
		}

		reqCont.LogHistory.Log("File successfully processed for upload")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(inserted_file.Filename + "." + inserted_file.Extension))
	}
}
