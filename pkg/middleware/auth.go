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

		authToken := r.Header.Get("Authorization")
		if authToken != "" {
			reqCont.IsLoggedIn = true
		}

		if globals.IsProtectedPath(r.URL.Path) && !reqCont.IsLoggedIn {
			http.Error(w, "Unauthorized: You must be logged in to access this resource", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
