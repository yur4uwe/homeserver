package middleware

import (
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"net/http"
)

func AuthMiddleware(connection *db.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
			if !ok {
				http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
				return
			}

			authToken, _ := r.Cookie("authTokenHS")
			if authToken == nil {
				reqCont.AuthToken = ""
			} else {
				_, err := connection.CheckToken(r.Context(), authToken.Value)
				if err == nil {
					reqCont.AuthToken = authToken.Value
				}
			}

			if globals.IsProtectedPath(r.URL.Path) && reqCont.AuthToken == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			reqCont.LogHistory = append(reqCont.LogHistory, "[AuthMiddleware] Request to "+r.URL.Path+" with token: "+authToken.Value)

			next.ServeHTTP(w, r)
		})
	}
}
