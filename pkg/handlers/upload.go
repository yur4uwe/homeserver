package handlers

import (
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/schemas"
	"homeserver/pkg/utils"
	"net/http"
	"path/filepath"
	"strings"
)

func downloadFile(connection *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
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

		file, header, err := r.FormFile("file")
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
		defer file.Close()

		path_str := r.FormValue("folder")
		if path_str != "" {
			path_segments := strings.Split(path_str, "/")
			path_str = reqCont.Requester.Username + "/" + strings.Join(path_segments[:len(path_segments)-1], "/")
			file_sig = strings.Split(path_segments[len(path_segments)-1], ".")
		} else {
			file_sig = strings.Split(header.Filename, ".")
			path_str = reqCont.Requester.Username
		}

		new_file := &schemas.UserFile{
			Filename:    file_sig[0],
			Path:        path_str,
			Extension:   file_sig[1],
			Size:        header.Size,
			ContentType: header.Header.Get("Content-Type"),
			OwnerID:     reqCont.Requester.ID,
		}

		_, saved_name, err := connection.SaveFile(r.Context(), new_file, &file, reqCont.Requester)
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

func deleteFile(connection *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Delete file request received")

		filePath := r.URL.Query().Get("filename")
		if filePath == "" {
			reqCont.LogHistory.ErrLog("Missing filename parameter", "DeleteFile")
			utils.ErrorResponse(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		file := schemas.NewUserFileBuilder().
			WithPath(filePath).
			WithOwnerID(reqCont.Requester.ID).
			WithFilename(strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))).
			WithExtension(filepath.Ext(filePath)[1:])

		if _, err := connection.DeleteFile(r.Context(), file, reqCont.Requester.Username); err != nil {
			reqCont.LogHistory.ErrLog("Error deleting file: "+err.Error(), "DeleteFile")
			utils.ErrorResponse(w, "Error deleting file", http.StatusInternalServerError)
			return
		}

		// Perform the file deletion logic here
		reqCont.LogHistory.Log("File successfully deleted: " + filePath)
		w.WriteHeader(http.StatusOK)
	}
}

func patchFile(connection *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
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

func UploadHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			downloadFile(connection, reqCont)(w, r)
		case http.MethodDelete:
			deleteFile(connection, reqCont)(w, r)
		case http.MethodPatch:
			patchFile(connection, reqCont)(w, r)
		default:
			reqCont.LogHistory.ErrLog("Unsupported method for upload handler: '"+r.Method+"'", "UploadHandler")
			utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
