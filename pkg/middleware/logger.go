package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type LogHistory struct {
	Entries []string
}

func (lh *LogHistory) AddLog(entry string) {
	lh.Entries = append(lh.Entries, entry)
}

func (lh *LogHistory) UnloadLogs() {
	logs := "\n" + strings.Join(lh.Entries, "\n")
	if len(logs) == 0 {
		log.Println("No logs to unload.")
		return
	}
	log.Println(logs)

}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		reqCont.LogHistory.AddLog(fmt.Sprintf("Request Method: %s, URL: %s", r.Method, r.URL.Path))

		next.ServeHTTP(w, r)
		reqCont.LogHistory.UnloadLogs()
	})
}
