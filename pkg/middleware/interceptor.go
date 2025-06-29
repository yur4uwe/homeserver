package middleware

import (
	"bytes"
	"homeserver/pkg/globals"
	"net/http"
	"strconv"
)

type responseInterceptor struct {
	statusCode int
	body       *bytes.Buffer
	http.ResponseWriter
}

func InterceptorConstructor(w http.ResponseWriter) *responseInterceptor {
	return &responseInterceptor{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           new(bytes.Buffer),
	}
}

func (ri *responseInterceptor) WriteHeader(code int) {
	ri.statusCode = code
}

func (ri *responseInterceptor) Write(b []byte) (int, error) {
	return ri.body.Write(b)
}

func Interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			http.Error(w, "Internal Server Error: Missing request context", http.StatusInternalServerError)
			return
		}

		if !reqCont.ToRestorePage {
			next.ServeHTTP(w, r)
			return
		}

		interceptor := InterceptorConstructor(w)

		next.ServeHTTP(interceptor, r)

		modifiedContent, err := globals.MergePage(interceptor.body.String())
		if err != nil {
			http.Error(w, "Failed to merge page content", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(modifiedContent)))
		w.WriteHeader(interceptor.statusCode)
		if _, err := w.Write(modifiedContent); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})
}
