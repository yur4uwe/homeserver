package handlers

import (
	"fmt"
	"homeserver/pkg/db"
	"homeserver/pkg/handlers/actions"
	"homeserver/pkg/handlers/search"
	"homeserver/pkg/handlers/upload"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

func UploadHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			upload.DownloadFile(connection, reqCont)(w, r)
		case http.MethodDelete:
			upload.DeleteFile(connection, reqCont)(w, r)
		case http.MethodPatch:
			upload.PatchFile(connection, reqCont)(w, r)
		default:
			reqCont.LogHistory.ErrLog("Unsupported method for upload handler: '"+r.Method+"'", "UploadHandler")
			utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

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
			actions.DeleteFileHandler(w, r, connection, reqCont, relativePath, fileName, isFolder)
		case http.MethodPatch:
			actions.RenameFileHandler(w, r, connection, reqCont, relativePath, fileName, desiredName, isFolder)
		case http.MethodGet:
			actions.OpenFileHandler(w, r, connection, reqCont, relativePath, fileName)
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

func SearchHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		query := r.URL.Query().Get("query")
		directory_string, err := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/search/"))
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API url.QueryUnescape error: "+err.Error(), "SearchHandler:30")
			utils.ErrorResponse(w, "Internal Server Error: Failed to unescape directory", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Search API accessed: query: '" + query + "', directory: '" + directory_string + "'")

		directory := []string{reqCont.Requester.Username}
		if directory_string != "" {
			directory = append(directory, strings.Split(directory_string, "/")...)
		}

		var data *templates.SearchPageData = nil

		if len(query) > 0 {
			data, err = search.FileSearchHandler(r.Context(), connection, query, directory, reqCont.Requester)
		} else {
			data, err = search.BrowseHandler(r.Context(), connection, directory, reqCont.Requester)
		}

		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			utils.ErrorResponse(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
			return
		}

		var parent []string
		if len(directory[1:]) > 0 {
			parent = directory[1 : len(directory)-1]
		}

		data.Parent = strings.Join(parent, "/")
		data.URIParent = search.ParseQuerySlice(parent)

		fmt.Println("SearchHandler: data:", data)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		var tmpl *template.Template
		if reqCont.ToRestorePage || reqCont.RerenderBody || r.URL.Query().Get("redirect") != "" {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_base.html", "../static/templates/search_results.html"))

			w.Header().Set("HX-Push-Url", "/search/"+url.QueryEscape(strings.Join(directory[1:], "/")))
			tmpl.ExecuteTemplate(w, "search_base.html", &data)
		} else {
			tmpl = template.Must(template.ParseFiles("../static/templates/search_results.html"))

			tmpl.Execute(w, data)
		}
	}
}
