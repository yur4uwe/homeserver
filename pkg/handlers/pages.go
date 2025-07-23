package handlers

import (
	"html/template"
	"net/http"

	"homeserver/pkg/middleware"
	"homeserver/pkg/templates"
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
