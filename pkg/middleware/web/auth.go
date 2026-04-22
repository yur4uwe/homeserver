package middleware

import (
	"fmt"
	"homeserver/constants"
	"homeserver/db"
	"homeserver/internal/templates"
	"homeserver/pkg/crypto"
	"html/template"
	"net/http"
	"strings"
)

var (
	ValidateToken = crypto.ValidateToken
)

var (
	headerTemplate *template.Template
)

func AuthMiddleware(repos *db.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/static/") {
				next.ServeHTTP(w, r)
				return
			}

			// Pre-parse template if not already done
			if headerTemplate == nil {
				var err error
				headerTemplate, err = template.ParseFiles(constants.Rootd + "static/templates/header.html")
				if err != nil {
					fmt.Println("CRITICAL: Failed to parse header template:", err)
				}
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

			claims, valid := ValidateToken(reqCont.AuthToken)
			if !valid || claims == nil {
				reqCont.AuthToken = ""
			}

			if user, err := repos.Tokens.QueryToken(r.Context(), reqCont.AuthToken, true); err == nil && user != nil && claims.Seed == user.Email {
				reqCont.Requester = user
				reqCont.LogHistory.Log(fmt.Sprintf("[AuthMiddleware] Authenticated user: %s", user.Username))
			} else {
				reqCont.AuthToken = ""
			}

			if IsProtectedPath(r.URL.Path) && reqCont.AuthToken == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if headerTemplate != nil {
				var buf strings.Builder
				if err := headerTemplate.Execute(&buf, templates.WelcomePageData{IsLoggedIn: reqCont.AuthToken != ""}); err != nil {
					reqCont.LogHistory.ErrLog("Failed to execute header template: "+err.Error(), "AuthMiddleware")
				} else {
					reqCont.PageContent.Header = template.HTML(buf.String())
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
