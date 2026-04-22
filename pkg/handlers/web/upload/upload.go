package upload

import (
	"homeserver/constants"
	"homeserver/db"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	Copy          = io.Copy
	MkdirAll      = os.MkdirAll
	Create        = os.Create
	Getenv        = os.Getenv
	ErrorResponse = utils.ErrorResponse
)

func DownloadFile(repos *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
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

		directory_id := 0

		path_parts := strings.Split(path_str, "/")

		for i, part := range path_parts {
			if part == "" {
				continue
			}

			curr_path := strings.Join(path_parts[:i], "/")

			fol := &db.UserFolder{
				Name:    part,
				OwnerID: reqCont.Requester.ID,
			}

			if curr_path != "" {
				fol.Path = curr_path
				fol.Parent = &directory_id
			}

			nfol, err := repos.Folders.Create(r.Context(), fol, reqCont.Requester)
			if err != nil {
				select {
				case <-done:
					// Handle cancellation
					reqCont.LogHistory.ErrLog("Request canceled during folder creation", "DownloadFile:66")
					ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
					return
				default:
					// Handle other errors
					reqCont.LogHistory.ErrLog("Error creating folder: "+err.Error(), "DownloadFile:71")
					ErrorResponse(w, "Error creating folder", http.StatusInternalServerError)
					return
				}
			}

			directory_id = nfol.ID
		}

		new_file := &db.UserFile{
			Filename:    file_sig[0],
			Path:        path_str,
			Extension:   file_sig[1],
			Size:        header.Size,
			Directory:   directory_id,
			ContentType: header.Header.Get("Content-Type"),
			OwnerID:     reqCont.Requester.ID,
		}

		inserted_file, err := repos.Files.Create(r.Context(), new_file, reqCont.Requester)
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

		if err := MkdirAll(Getenv(constants.StoragePath)+inserted_file.Path, os.ModePerm); err != nil {
			select {
			case <-done:
				// Handle cancellation
				reqCont.LogHistory.ErrLog("Request canceled during file saving", "DownloadFile:66")
				ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
				return
			default:
				// Handle other errors
				reqCont.LogHistory.ErrLog("Error creating directories: "+err.Error(), "DownloadFile:71")
				ErrorResponse(w, "Error creating directories", http.StatusInternalServerError)
				return
			}
		}

		f, err := Create(inserted_file.AbsolutePath())
		if err != nil {
			select {
			case <-done:
				// Handle cancellation
				reqCont.LogHistory.ErrLog("Request canceled during file saving", "DownloadFile:66")
				ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
				return
			default:
				// Handle other errors
				reqCont.LogHistory.ErrLog("Error creating file: "+err.Error(), "DownloadFile:71")
				ErrorResponse(w, "Error creating file", http.StatusInternalServerError)
				return
			}
		}
		defer f.Close()

		if _, err := Copy(f, multipart); err != nil {
			select {
			case <-done:
				// Handle cancellation
				reqCont.LogHistory.ErrLog("Request canceled during file saving", "DownloadFile:66")
				ErrorResponse(w, "Request canceled", http.StatusRequestTimeout)
				return
			default:
				// Handle other errors
				reqCont.LogHistory.ErrLog("Error copying file content: "+err.Error(), "DownloadFile:71")
				ErrorResponse(w, "Error copying file content", http.StatusInternalServerError)
				return
			}
		}

		reqCont.LogHistory.Log("File successfully processed for upload")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(inserted_file.Filename + "." + inserted_file.Extension))
	}
}
