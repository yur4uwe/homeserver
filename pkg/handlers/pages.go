package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/templates"

	"golang.org/x/net/html"
)

func WelcomePageHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Welcome page accessed")

	tmpl := template.Must(template.ParseFiles("../static/templates/welcome.html"))

	data := templates.WelcomePageData{
		IsLoggedIn: reqCont.AuthToken != "",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Upload page accessed")

	fragment := parser.Must(parser.ParseFile("../static/fragments/upload.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html.Render(w, fragment)
}

func SearchPageHandler(data *templates.SearchPageData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("[Handler] Search page accessed")

		tmpl := template.Must(template.ParseFiles("../static/templates/search.html"))

		directory := strings.TrimPrefix(r.URL.Path, "/search/")

		reqCont.LogHistory.Log("[Handler] Search directory: " + directory)

		if data == nil {
			data = &templates.SearchPageData{
				Query:     "",
				Directory: directory,
				Entries:   nil,
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}
