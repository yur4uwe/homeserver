package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func AddLog(r *http.Request, entry string) error {
	reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
	if !ok {
		return fmt.Errorf("request context not found")
	}

	reqCont.LogHistory = append(reqCont.LogHistory, entry)
	return nil
}

func UnloadLogs(r *http.Request) {
	finalCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
	if !ok {
		log.Println("No request context found to unload logs.")
		return
	}

	logs := "\n" + strings.Join(finalCont.LogHistory, "\n")
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

		reqCont.LogHistory = append(reqCont.LogHistory, fmt.Sprintf("[Logger] Request to %s %s", r.Method, r.URL.Path))

		next.ServeHTTP(w, r)
		UnloadLogs(r)
	})
}
