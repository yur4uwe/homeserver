package middleware

import (
	"context"
	"homeserver/pkg/globals"
	"net/http"
)

type contextKey string

const ReqContKey contextKey = "add_req_context"

type AdditionalRequestContext struct {
	LogHistory    []string
	ToRestorePage bool
}

func contextConstructor() *AdditionalRequestContext {
	return &AdditionalRequestContext{
		LogHistory:    []string{},
		ToRestorePage: false,
	}
}

func ToRestorePage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") != "true" && globals.IsPage(r.URL.Path) {
			reqCont.ToRestorePage = true
			err := AddLog(r, "[Middleware] ToRestorePage set to true")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func WithAdditionalRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont := contextConstructor()
		ctx := context.WithValue(r.Context(), ReqContKey, reqCont)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func MiddlewareWrapper(handler http.Handler) http.Handler {
	return WithAdditionalRequestContext(Logger(ToRestorePage(Interceptor(handler))))
}
