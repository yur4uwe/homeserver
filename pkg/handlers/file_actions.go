package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func DeleteFileHandler(connection *db.Database, relative_path []string, file_name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		file_path := strings.Join(append(relative_path, file_name), "/")

		err := connection.DeleteFile(r.Context(), file_path, reqCont.Requester)
		if err != nil {
			utils.ErrorResponse(w, fmt.Sprintf("Error deleting file: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func RenameFileHandler(connection *db.Database, relative_path []string, file_name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		new_name := r.URL.Query().Get("new_name")
		if new_name == "" {
			utils.ErrorResponse(w, "New name parameter is required", http.StatusBadRequest)
			return
		}

		file_path := strings.Join(append(relative_path, file_name), "/")

		actual_name, err := connection.RenameFile(r.Context(), file_path, new_name, reqCont.Requester)
		if err != nil {
			utils.ErrorResponse(w, fmt.Sprintf("Error renaming file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(actual_name))
	}
}

func OpenFileHandler(connection *db.Database, relative_path []string, file_name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		file_path := strings.Join(append(relative_path, file_name), "/")

		file, err := connection.GetFile(r.Context(), file_path, reqCont.Requester)
		if err != nil {
			utils.ErrorResponse(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, file.Path)
	}
}

func FileActionsHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		action := r.URL.Query().Get("action")
		if action == "" {
			utils.ErrorResponse(w, "Action parameter is required", http.StatusBadRequest)
			return
		}
		is_folder := r.URL.Query().Get("is_folder") == "true"

		path := append([]string{reqCont.Requester.Username}, strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/file-actions/"), "/")...)
		relativePath := path[:len(path)-1]
		file_name := path[len(path)-1]

		reqCont.LogHistory.Log(fmt.Sprintf("FileActionsHandler called: action=%s, is_folder=%t, relativePath=%v, file_name=%s", action, is_folder, relativePath, file_name))

		switch action {
		case "delete":
			DeleteFileHandler(connection, relativePath, file_name)(w, r)
		case "rename":
			RenameFileHandler(connection, relativePath, file_name)(w, r)
		case "open":
			OpenFileHandler(connection, relativePath, file_name)(w, r)
		default:
			utils.ErrorResponse(w, fmt.Sprintf("Unknown action: %s", action), http.StatusBadRequest)
			return
		}

		alert := parser.NewAlert(
			"File Action Success",
			"check",
			fmt.Sprintf("Action '%s' performed on file '%s'.", action, file_name),
			utils.AlertSuccess,
		)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("HX-Retarget", "#modal-content")
		w.Header().Set("HX-Reswap", "innerHTML")
		w.Write([]byte(alert))
	}
}
