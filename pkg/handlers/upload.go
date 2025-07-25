package handlers

import (
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"net/http"
	"strings"
)

func downloadFile(connection *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Download file request received")

		var file_name, relative_path string

		file, header, err := r.FormFile("file")
		if err != nil {
			reqCont.LogHistory.ErrLog("Error retrieving file from form: "+err.Error(), "DownloadFile")
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		relative_path = r.FormValue("folder")
		if relative_path != "" {
			path_segments := strings.Split(relative_path, "/")

			relative_path = strings.Join(path_segments[:len(path_segments)-1], "/")

			file_name = path_segments[len(path_segments)-1]
		} else {
			file_name = header.Filename
			relative_path = ""
		}

		fileInfo := &templates.FileInfo{
			Name:        file_name,
			Path:        relative_path,
			Size:        header.Size,
			ContentType: header.Header.Get("Content-Type"),
		}

		_, saved_name, err := connection.SaveFile(r.Context(), fileInfo, &file, reqCont.AuthToken)
		if err != nil {
			reqCont.LogHistory.ErrLog("Error saving file: "+err.Error(), "DownloadFile")
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
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
			http.Error(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		if err := connection.DeleteFile(r.Context(), filePath, reqCont.AuthToken); err != nil {
			reqCont.LogHistory.ErrLog("Error deleting file: "+err.Error(), "DeleteFile")
			http.Error(w, "Error deleting file", http.StatusInternalServerError)
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

		newFilePath, err := connection.RenameFile(r.Context(), oldFileName, newFileName, reqCont.AuthToken)
		if err != nil {
			reqCont.LogHistory.ErrLog("Error renaming file: "+err.Error(), "PatchFile")
			http.Error(w, "Error renaming file", http.StatusInternalServerError)
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
			http.Error(w, "Invalid request context", http.StatusInternalServerError)
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
