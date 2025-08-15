package middleware

import (
	"fmt"
	"homeserver/pkg/utils"
	"net/http"
	"strings"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := Now()

		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			utils.ErrorResponse(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		query := r.URL.RawQuery
		url := r.URL.Path
		if query != "" {
			url += "?" + query
		}

		next.ServeHTTP(w, r)

		duration := Since(start)

		header := fmt.Sprintf("Request: %s %s took %v\n", r.Method, url, duration)

		if !strings.HasPrefix(url, "/static/") {
			reqCont.LogHistory.UnloadLogs(header)
		}
	})
}
