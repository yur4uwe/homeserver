package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"homeserver/pkg/globals"
	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
)

func WelcomePageHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
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
		utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Upload page accessed")

	tmpl := template.Must(template.ParseFiles("../static/fragments/upload.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func StaticPageHandler(w http.ResponseWriter, r *http.Request) {
	reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
	if !ok {
		fmt.Fprintf(os.Stderr, "Missing request context in StaticPageHandler\n")
		utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
		return
	}

	reqCont.LogHistory.Log("[Handler] Page server accessed for page: " + r.URL.Path)

	if !globals.IsPage(r.URL.Path) {
		reqCont.LogHistory.Log("Path not considered a page: " + r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/" {
		WelcomePageHandler(w, r)
		return
	}

	file_to_serve := "../static/fragments/" + r.URL.Path[1:] + ".html"

	content, err := os.ReadFile(file_to_serve)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file %s: %v\n", file_to_serve, err.Error())
		utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write reponse: %v\n", err.Error())
		utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
