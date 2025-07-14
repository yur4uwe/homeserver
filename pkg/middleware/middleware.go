package middleware

import (
	"context"
	"homeserver/pkg/db"
	"homeserver/pkg/globals"
	"net/http"

	"golang.org/x/net/html"
)

type contextKey string

const ReqContKey contextKey = "add_req_context"

type AdditionalRequestContext struct {
	LogHistory     *LogHistory
	ToRestorePage  bool
	FullPageRender *html.Node
	AuthToken      string
}

func contextConstructor() *AdditionalRequestContext {
	return &AdditionalRequestContext{
		LogHistory:     &LogHistory{Entries: []string{}},
		ToRestorePage:  false,
		FullPageRender: nil,
		AuthToken:      "",
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
			reqCont.LogHistory.Log("[Middleware] ToRestorePage set to true")
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

func MiddlewareWrapper(handler http.Handler, connection *db.Database) http.Handler {
	return WithAdditionalRequestContext(Logger(AuthMiddleware(connection)(ToRestorePage(Interceptor(handler)))))
}
