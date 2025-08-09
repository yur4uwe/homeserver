package auth

import (
	"homeserver/pkg/interfaces"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/templates"
	"html/template"
	"net/http"
)

var (
	ParseTemplateToString = parser.ParseTemplateToString
)

func LogoutHandler(connection interfaces.DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("Logout request received")

		authToken := reqCont.AuthToken
		if authToken == "" {
			ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := connection.DeleteToken(r.Context(), authToken); err != nil {
			ErrorResponse(w, "Failed to delete token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reqCont.RerenderBody = true

		reqCont.AuthToken = ""
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("HX-Push-Url", "/")

		data := templates.WelcomePageData{
			IsLoggedIn: false,
		}

		header, err := ParseTemplateToString("../static/templates/header.html", data)
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse header template: "+err.Error(), "LogoutHandler:66")
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reqCont.PageContent.Header = template.HTML(header)

		tmpl := Must(ParseFiles("../static/templates/welcome.html"))
		if err := tmpl.Execute(w, data); err != nil {
			reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "LogoutHandler:66")
			ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
