package middleware

import (
	"fmt"
	"homeserver/pkg/interfaces"
	"homeserver/pkg/templates"
	"html/template"
	"net/http"
	"strings"
)

func AuthMiddleware(connection interfaces.DatabaseInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/static/") {
				next.ServeHTTP(w, r)
				return
			}

			reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
			if !ok {
				ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
				return
			}

			authToken, _ := r.Cookie("authTokenHS")
			if authToken == nil {
				reqCont.AuthToken = ""
			} else {
				reqCont.AuthToken = authToken.Value
			}

			if user, err := connection.QueryToken(r.Context(), reqCont.AuthToken); err == nil && user != nil {
				reqCont.Requester = user
			} else {
				reqCont.AuthToken = ""
			}

			if IsProtectedPath(r.URL.Path) && reqCont.AuthToken == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if header, err := ParseTemplateToString("../static/templates/header.html", templates.WelcomePageData{IsLoggedIn: reqCont.AuthToken != ""}); err != nil {
				fmt.Println("Failed to parse header template:", err)
				reqCont.LogHistory.ErrLog("Failed to parse header template: "+err.Error(), "AuthMiddleware")
				ErrorResponse(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
				return
			} else {
				fmt.Println("Parsed header template successfully:", header)
				reqCont.PageContent.Header = template.HTML(header)
			}

			reqCont.LogHistory.Log("[AuthMiddleware] Request to " + r.URL.Path + " with token: '" + reqCont.AuthToken + "'")
			next.ServeHTTP(w, r)
		})
	}
}
