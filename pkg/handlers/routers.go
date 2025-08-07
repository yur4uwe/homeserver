package handlers

import (
	"fmt"
	"homeserver/pkg/handlers/actions"
	"homeserver/pkg/handlers/search"
	"homeserver/pkg/handlers/upload"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrorResponse      = utils.ErrorResponse
	ParsePathWithUName = utils.ParsePathWithUName

	DownloadFile = upload.DownloadFile
	DeleteFile   = upload.DeleteFile
	PatchFile    = upload.PatchFile

	DeleteFileHandler = actions.DeleteFileHandler
	RenameFileHandler = actions.RenameFileHandler
	OpenFileHandler   = actions.OpenFileHandler

	FileSearchHandler = search.FileSearchHandler
	BrowseHandler     = search.BrowseHandler
)

func UploadRouter(connection interfaces.DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			DownloadFile(connection, reqCont)(w, r)
		case http.MethodDelete:
			DeleteFile(connection, reqCont)(w, r)
		case http.MethodPatch:
			PatchFile(connection, reqCont)(w, r)
		default:
			reqCont.LogHistory.ErrLog("Unsupported method for upload handler: '"+r.Method+"'", "UploadHandler")
			ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func FileActionsRouter(connection interfaces.DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		isFolder := r.URL.Query().Get("is_folder") == "true"
		desiredName := r.URL.Query().Get("desired_name")

		entryPath := strings.TrimPrefix(r.URL.Path, "/api/v1/file-actions/")
		if entryPath == "" {
			ErrorResponse(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		entryName, pathSlice, err := ParsePathWithUName(entryPath, reqCont.Requester.Username)
		if err != nil {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqCont.LogHistory.Log(fmt.Sprintf("FileActionsRouter called: method=%s, is_folder=%t, relativePath=%v, file_name=%s", r.Method, isFolder, pathSlice, entryName))

		switch r.Method {
		case http.MethodDelete:
			DeleteFileHandler(w, r, connection, reqCont, pathSlice, entryName, isFolder)
		case http.MethodPatch:
			RenameFileHandler(w, r, connection, reqCont, pathSlice, entryName, desiredName, isFolder)
		case http.MethodGet:
			OpenFileHandler(w, r, connection, reqCont, pathSlice, entryName)
		default:
			ErrorResponse(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

func SearchRouter(connection interfaces.DatabaseInterface) http.HandlerFunc {
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
			data, err = FileSearchHandler(r.Context(), connection, query, directory, reqCont.Requester)
		} else {
			data, err = BrowseHandler(r.Context(), connection, directory, reqCont.Requester)
		}

		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API db query error: "+err.Error(), "SearchHandler:24")
			ErrorResponse(w, "Internal Server Error: Failed to search files", http.StatusInternalServerError)
			return
		}

		var parent []string
		if len(directory[1:]) > 0 {
			parent = directory[1 : len(directory)-1]
		}

		data.Parent = strings.Join(parent, "/")
		data.URIParent = search.ParseQuerySlice(parent)

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
