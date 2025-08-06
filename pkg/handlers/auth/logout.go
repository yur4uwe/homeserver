package auth

import (
	"homeserver/pkg/db"
	"homeserver/pkg/middleware"
	"homeserver/pkg/parser"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
)

func LogoutHandler(connection *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(middleware.ReqContKey).(*middleware.AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Invalid request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.Log("Logout request received")

		authToken := reqCont.AuthToken
		if authToken == "" {
			utils.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := connection.DeleteToken(r.Context(), authToken); err != nil {
			utils.ErrorResponse(w, "Failed to delete token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reqCont.RerenderBody = true

		reqCont.AuthToken = ""
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("HX-Push-Url", "/")
		w.WriteHeader(http.StatusOK)

		data := templates.WelcomePageData{
			IsLoggedIn: false,
		}

		header, err := parser.ParseTemplateToString("../static/templates/header.html", data)
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse header template: "+err.Error(), "LogoutHandler:66")
			utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reqCont.PageContent.Header = template.HTML(header)

		tmpl := template.Must(template.ParseFiles("../static/templates/welcome.html"))
		if err := tmpl.Execute(w, data); err != nil {
			reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "LogoutHandler:66")
			utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
