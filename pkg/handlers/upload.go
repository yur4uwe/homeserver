package handlers

import (
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"net/http"
	"strconv"
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

		fileInfo := templates.FileInfo{
			Name:        file_name,
			Path:        relative_path,
			Size:        header.Size,
			ContentType: header.Header.Get("Content-Type"),
		}

		reqCont.LogHistory.Log("File received: " + fileInfo.Name)
		reqCont.LogHistory.Log("Relative path: '" + fileInfo.Path + "'")
		reqCont.LogHistory.Log("File size: " + strconv.FormatInt(fileInfo.Size, 10) + " bytes")
		reqCont.LogHistory.Log("File content type: " + fileInfo.ContentType)

		// connection.SaveFile(r.Context(), &fileInfo, &file, "/", reqCont.AuthToken)

		// Process the file (e.g., save it to the server)

		reqCont.LogHistory.Log("File successfully processed for upload")
		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteFile(connection *db.Database, reqCont *middleware.AdditionalRequestContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont.LogHistory.Log("Delete file request received")

		fileName := r.URL.Query().Get("filename")
		if fileName == "" {
			reqCont.LogHistory.ErrLog("Missing filename parameter", "DeleteFile")
			http.Error(w, "Missing filename parameter", http.StatusBadRequest)
			return
		}

		// Perform the file deletion logic here
		reqCont.LogHistory.Log("File successfully deleted: " + fileName)
		w.WriteHeader(http.StatusNoContent)
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
		default:
			reqCont.LogHistory.ErrLog("Unsupported method for upload handler: '"+r.Method+"'", "UploadHandler")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
