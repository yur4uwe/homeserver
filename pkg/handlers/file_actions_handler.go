package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func FileActionsRouter(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		isFolder := r.URL.Query().Get("is_folder") == "true"
		desiredName := r.URL.Query().Get("desired_name")

		path := append([]string{reqCont.Requester.Username}, strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/file-actions/"), "/")...)
		relativePath := path[:len(path)-1]
		fileName := path[len(path)-1]

		reqCont.LogHistory.Log(fmt.Sprintf("FileActionsRouter called: method=%s, is_folder=%t, relativePath=%v, file_name=%s", r.Method, isFolder, relativePath, fileName))

		switch r.Method {
		case http.MethodDelete:
			DeleteFileHandler(w, r, connection, reqCont, relativePath, fileName, isFolder)
		case http.MethodPatch:
			RenameFileHandler(w, r, connection, reqCont, relativePath, fileName, desiredName, isFolder)
		case http.MethodGet:
			OpenFileHandler(w, r, connection, reqCont, relativePath, fileName)
		default:
			utils.ErrorResponse(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}

		// alert := parser.NewAlert(
		// 	"File Action Success",
		// 	"check",
		// 	fmt.Sprintf("Action '%s' performed on file '%s'.", action, fileName),
		// 	utils.AlertSuccess,
		// )

		// w.Header().Set("HX-Trigger", fmt.Sprintf(`%sSuccess-%s-%s`, action, strings.Join(relativePath, "/"), fileName))
		// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Header().Set("HX-Retarget", "#modal-content")
		// w.Header().Set("HX-Reswap", "innerHTML")
		// w.Write([]byte(alert))
	}
}
