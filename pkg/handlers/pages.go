package handlers

import (
	"fmt"
	"net/http"
	"os"

	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
)

var (
	ReadFile = os.ReadFile

	IsPage = globals.IsPage

	WelcomePage = WelcomePageHandler
)

func WelcomePageHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Welcome page accessed")

	tmpl := Must(ParseFiles("../static/templates/welcome.html"))

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
		ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Upload page accessed")

	tmpl := Must(ParseFiles("../static/fragments/upload.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func SettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Settings page accessed")

	tmpl := Must(ParseFiles("../static/templates/settings.html"))

	data := templates.SettingsPageData{
		Username: reqCont.Requester.Username,
		Email:    reqCont.Requester.Email,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, data); err != nil {
		ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	}
}

func StaticPageHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		fmt.Fprintf(os.Stderr, "Missing request context in StaticPageHandler\n")
		ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Page server accessed for page: " + r.URL.Path)

	if !IsPage(r.URL.Path) {
		reqCont.LogHistory.Log("Path not considered a page: " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/" {
		WelcomePage(w, r)
		return
	}

	file_to_serve := "../static/fragments/" + r.URL.Path[1:] + ".html"

	content, err := ReadFile(file_to_serve)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file %s: %v\n", file_to_serve, err.Error())
		ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}
