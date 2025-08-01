package middleware

import (
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"homeserver/pkg/parser"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
	"html/template"
	"net/http"
	"time"
)

func AuthMiddleware(connection *db.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
			if !ok {
				utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
				return
			}

			authToken, _ := r.Cookie("authTokenHS")
			if authToken == nil {
				reqCont.AuthToken = ""
			} else if authToken.Expires.After(time.Now()) {
				reqCont.AuthToken = ""

				http.SetCookie(w, &http.Cookie{
					Name:     "authTokenHS",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
			} else {
				reqCont.AuthToken = authToken.Value
			}

			user, err := connection.QueryToken(r.Context(), reqCont.AuthToken)
			if err == nil && user != nil {
				reqCont.Requester = user
			} else {
				reqCont.AuthToken = ""
			}

			if globals.IsProtectedPath(r.URL.Path) && reqCont.AuthToken == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			data := templates.WelcomePageData{
				IsLoggedIn: reqCont.AuthToken != "",
			}

			header, err := parser.ParseTemplateToString("../static/templates/header.html", data)
			if err != nil {
				reqCont.LogHistory.ErrLog("Failed to parse header template: "+err.Error(), "AuthMiddleware:66")
				utils.ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			reqCont.PageContent.Header = template.HTML(header)

			reqCont.LogHistory.Log("[AuthMiddleware] Request to " + r.URL.Path + " with token: '" + reqCont.AuthToken + "'")

			next.ServeHTTP(w, r)
		})
	}
}
