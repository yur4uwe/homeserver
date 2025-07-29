package handlers

import (
	"html/template"
	"net/http"

	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"

	"golang.org/x/net/html"
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

	fragment := parser.Must(parser.ParseFile("../static/fragments/upload.html"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html.Render(w, fragment)
}
