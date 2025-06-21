package middleware

import (
	"fmt"
	"log"
	"net/http"
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

	for _, entry := range finalCont.LogHistory {
		log.Println(entry)
	}
}
