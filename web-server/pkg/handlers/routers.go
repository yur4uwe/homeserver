package handlers

import (
	"fmt"
	"homeserver/db/repositories"
	"homeserver/web-server/pkg/handlers/account"
	"homeserver/web-server/pkg/handlers/actions"
	"homeserver/web-server/pkg/handlers/search"
	"homeserver/web-server/pkg/handlers/upload"
	"homeserver/web-server/pkg/middleware"
	"homeserver/web-server/pkg/templates"
	"homeserver/web-server/pkg/utils"
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

	FileSearchHandler = search.FileSearchHandler
	BrowseHandler     = search.BrowseHandler
	ParseQuerySlice   = search.ParseQuerySlice

	Must       = template.Must
	ParseFiles = template.ParseFiles
)

func UploadRouter(repos *repositories.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			DownloadFile(repos, reqCont)(w, r)
		case http.MethodDelete:
			DeleteFile(repos, reqCont)(w, r)
		case http.MethodPatch:
			PatchFile(repos, reqCont)(w, r)
		default:
			reqCont.LogHistory.ErrLog("Unsupported method for upload handler: '"+r.Method+"'", "UploadHandler")
			ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func FileActionsRouter(repos *repositories.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		isFolder := r.URL.Query().Get("is_folder") == "true"
		desiredName := r.URL.Query().Get("desired_name")

		entryPath, err := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/api/v1/file-actions/"))
		if err != nil || entryPath == "" {
			ErrorResponse(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		pathSlice, err := ParsePathWithUName(entryPath, reqCont.Requester.Username)
		if err != nil {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		relativePath := pathSlice[:len(pathSlice)-1]
		fileName := pathSlice[len(pathSlice)-1]

		reqCont.LogHistory.Log(
			fmt.Sprintf(
				"FileActionsRouter called: method=%s, is_folder=%t, relativePath=%v, file_name=%s",
				r.Method, isFolder, relativePath, fileName,
			),
		)

		switch r.Method {
		case http.MethodDelete:
			DeleteFileHandler(w, r, repos, reqCont, relativePath, fileName)
		case http.MethodPatch:
			RenameFileHandler(w, r, repos, reqCont, relativePath, fileName, desiredName, isFolder)
		default:
			ErrorResponse(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

func SearchRouter(repos *repositories.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		query := r.URL.Query().Get("query")
		directory_string, err := url.QueryUnescape(strings.TrimPrefix(r.URL.Path, "/search/"))
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API url.QueryUnescape error: "+err.Error(), "SearchHandler:30")
			ErrorResponse(w, "Internal Server Error: Failed to unescape directory", http.StatusBadRequest)
			return
		}

		reqCont.LogHistory.Log("[Handler] Search API accessed: query: '" + query + "', directory: '" + directory_string + "'")

		directory, err := ParsePathWithUName(directory_string, reqCont.Requester.Username)
		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Search API ParsePathWithUName error: "+err.Error(), "SearchHandler:20")
			ErrorResponse(w, "Internal Server Error: Failed to parse directory", http.StatusInternalServerError)
			return
		}

		var data *templates.SearchPageData = nil

		if len(query) > 0 {
			data, err = FileSearchHandler(r.Context(), repos, query, directory, reqCont.Requester)
		} else {
			data, err = BrowseHandler(r.Context(), repos, directory, reqCont.Requester)
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
		data.URIParent = ParseQuerySlice(parent)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		var tmpl *template.Template
		if reqCont.ToRestorePage || reqCont.RerenderBody || r.URL.Query().Get("redirect") != "" {
			tmpl = Must(ParseFiles("../static/templates/search_base.html", "../static/templates/search_results.html"))

			w.Header().Set("HX-Push-Url", "/search/"+url.QueryEscape(strings.Join(directory[1:], "/")))
			tmpl.ExecuteTemplate(w, "search_base.html", &data)
		} else {
			tmpl = Must(ParseFiles("../static/templates/search_results.html"))

			tmpl.Execute(w, data)
		}
	}
}

func AccountSettingsRouter(repos *repositories.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Account settings accessed")

		new_username := r.FormValue("username")
		new_email := r.FormValue("email")
		new_password := r.FormValue("password")

		var err error
		if new_username != "" && new_username != reqCont.Requester.Username {
			reqCont.LogHistory.Log("[Handler] Username change requested")
			if uname_err := account.HandleUsernameChange(r, repos); uname_err != nil {
				err = uname_err
			}
		}

		if new_email != "" && new_email != reqCont.Requester.Email {
			reqCont.LogHistory.Log("[Handler] Email change requested")
			if em_err := account.HandleEmailChange(r, repos); em_err != nil {
				err = em_err
			}
		}

		if new_password != "" && new_password == r.FormValue("confirm-password") {
			reqCont.LogHistory.Log("[Handler] Password change requested")
			if pass_err := account.HandlePasswordChange(r, repos); pass_err != nil {
				err = pass_err
			}
		}

		if err != nil {
			reqCont.LogHistory.ErrLog("[Handler] Account settings update error: "+err.Error(), "AccountSettingsRouter:40")
			ErrorResponse(w, "Internal Server Error: Failed to update account settings", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Account settings updated successfully"))
	}
}
