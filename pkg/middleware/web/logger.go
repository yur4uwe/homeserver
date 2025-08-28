package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := Now()

		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		query := r.URL.RawQuery
		url := r.URL.Path
		if query != "" {
			url += "?" + query
		}

		next.ServeHTTP(w, r)

		duration := Since(start)

		if !strings.HasPrefix(url, "/static/") {
			header := fmt.Sprintf("Request: %s %s took %v\n", r.Method, url, duration)
			reqCont.LogHistory.UnloadLogs(header)
		}
	})
}
