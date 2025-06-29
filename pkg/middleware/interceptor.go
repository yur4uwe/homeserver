package middleware

import (
	"bytes"
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
	// Delay writing header until after handler finishes
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
			// If the request is to restore a page, we skip the interceptor
			next.ServeHTTP(w, r)
			return
		}
		// Create a response interceptor
		interceptor := InterceptorConstructor(w)

		// Call the next handler with the interceptor
		next.ServeHTTP(interceptor, r)

		modifiedContent := []byte(interceptor.body.String() + "\n<!-- Intercepted Response -->")

		w.Header().Set("Content-Length", strconv.Itoa(len(modifiedContent)))
		w.WriteHeader(interceptor.statusCode)
		if _, err := w.Write(modifiedContent); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})
}
