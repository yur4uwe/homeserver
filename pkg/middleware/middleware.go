package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const ReqContKey contextKey = "add_req_context"

type AdditionalRequestContext struct {
	LogHistory []string
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont := &AdditionalRequestContext{
			LogHistory: []string{"[Initial] Request received: " + r.Method + " " + r.URL.String()},
		}

		ctx := context.WithValue(r.Context(), ReqContKey, reqCont)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		UnloadLogs(r)
	})
}

func MiddlewareWrapper(handler http.Handler) http.Handler {
	return Logger(handler)
}
