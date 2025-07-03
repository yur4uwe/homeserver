package middleware

import (
	"homeserver/pkg/globals"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		authToken, err := r.Cookie("authTokenHS")
		if err == nil && authToken != nil {
			reqCont.IsLoggedIn = true
		}

		if globals.IsProtectedPath(r.URL.Path) && !reqCont.IsLoggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		reqCont.LogHistory = append(reqCont.LogHistory, "[AuthMiddleware] Request to "+r.URL.Path+" with token: "+authToken.Value)

		next.ServeHTTP(w, r)
	})
}
