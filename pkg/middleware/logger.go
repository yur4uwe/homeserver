package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type LogHistory struct {
	Entries []string
}

func (lh *LogHistory) Log(entry string) {
	lh.Entries = append(lh.Entries, entry)
}

func (lh *LogHistory) ErrLog(entry, at string) {
	lh.Entries = append(lh.Entries, "[Err at "+at+"] "+entry)
}

func (lh *LogHistory) UnloadLogs(header string) {
	logs := strings.Join(lh.Entries, "\n") + "\n\n"
	log.Println("\n\n" + header + logs)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		query := r.URL.RawQuery
		url := r.URL.Path
		if query != "" {
			url += "?" + query
		}

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		header := fmt.Sprintf("Request: %s %s took %v\n", r.Method, url, duration)

		if !strings.HasPrefix(url, "/static/") {
			reqCont.LogHistory.UnloadLogs(header)
		}
	})
}
