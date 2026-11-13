package handlers

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db"
	"homeserver/internal/templates"
	"homeserver/pkg/handlers/web/account"
	"homeserver/pkg/handlers/web/search"
	"homeserver/pkg/handlers/web/upload"
	middleware "homeserver/pkg/middleware/web"
	"homeserver/pkg/parser"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrorResponse      = utils.ErrorResponse
	ParsePathWithUName = parser.ParsePathWithUName

	DownloadFile = upload.DownloadFile
	DeleteFile   = upload.DeleteFile
	PatchFile    = upload.PatchFile

	FileSearchHandler = search.FileSearchHandler
	BrowseHandler     = search.BrowseHandler
	ParseQuerySlice   = search.ParseQuerySlice

	Must       = template.Must
	ParseFiles = template.ParseFiles
)

func UploadRouter(repos *db.Database) http.HandlerFunc {
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

func SearchRouter(repos *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		query := r.URL.Query().Get("query")
		directory_string, err := url.QueryUnescape(strings.TrimPrefix(r.URL.EscapedPath(), "/search/"))
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
			var page int64 = 1
			if raw := r.URL.Query().Get("page"); raw != "" {
				page, err = strconv.ParseInt(raw, 10, 32)
				if err != nil || page < 1 {
					ErrorResponse(w, "Invalid search page", http.StatusBadRequest)
					return
				}
			}
			data, err = FileSearchHandler(r.Context(), repos, query, directory, reqCont.Requester, int32(page))
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
		data.DownloadPath = fmt.Sprintf("http://%s:%s", Getenv(constants.FileSrvHost), Getenv(constants.FileSrvPort))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		var tmpl *template.Template
		if reqCont.ToRestorePage || reqCont.RerenderBody || r.URL.Query().Get("redirect") != "" {
			tmpl = Must(ParseFiles(constants.Rootd+"static/templates/search_base.html", constants.Rootd+"static/templates/search_results.html"))

			searchURL := "/search/" + data.URIDirectory
			if query != "" {
				values := url.Values{"query": {query}}
				if page := r.URL.Query().Get("page"); page != "" {
					values.Set("page", page)
				}
				searchURL += "?" + values.Encode()
			}
			w.Header().Set("HX-Push-Url", searchURL)
			tmpl.ExecuteTemplate(w, "search_base.html", &data)
		} else {
			tmpl = Must(ParseFiles(constants.Rootd + "static/templates/search_results.html"))

			tmpl.Execute(w, data)
		}
	}
}

func AccountSettingsRouter(repos *db.Database) http.HandlerFunc {
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

		alertMsg := ""
		data := templates.SettingsPageData{
			Username: reqCont.Requester.Username,
			Email:    reqCont.Requester.Email,
		}

		var err = map[string]string{}
		if uname_err := account.HandleUsernameChange(r.Context(), repos, reqCont, new_username); uname_err != nil {
			data.Username = reqCont.Requester.Username
			err["username"] = uname_err.Error()
		} else {
			data.Username = new_username
			alertMsg += "Username changed successfully. "
		}

		if em_err := account.HandleEmailChange(r.Context(), repos, reqCont, new_email); em_err != nil {
			data.Email = reqCont.Requester.Email
			err["email"] = em_err.Error()
		} else {
			data.Email = new_email
			alertMsg += "Email changed successfully. "
		}

		if pass_err := account.HandlePasswordChange(r.Context(), repos, reqCont, new_password, r.FormValue("confirm-password")); pass_err != nil {
			err["password"] = pass_err.Error()
		} else {
			alertMsg += "Password changed successfully. "
		}

		if len(err) == 3 {
			data.Username = reqCont.Requester.Username
			data.Email = reqCont.Requester.Email

			var errs strings.Builder
			for k, v := range err {
				fmt.Fprintf(&errs, "%s: %s\n", k, v)
			}

			data.Error = errs.String()
		} else if alertMsg != "" {
			w.Header().Set("HX-Trigger", utils.AlertEvent("Success!", alertMsg, constants.Success))
		}

		if strings.Contains(alertMsg, "Email") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tmpl := Must(ParseFiles(constants.Rootd + "static/templates/settings.html"))
		tmpl.Execute(w, data)
	}
}
