package middleware

import (
	"bytes"
	"fmt"
	"homeserver/pkg/parser"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/net/html"
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

		interceptor := InterceptorConstructor(w)

		next.ServeHTTP(interceptor, r)

		body_node, err := parser.ParseString(interceptor.body.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse response body: %v\n", err.Error())
			http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
			return
		}

		var modifiedContent *html.Node = nil

		if reqCont.ToRestorePage {
			modifiedContent, err = parser.MergePage(body_node)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to merge page content: %v\n", err.Error())
				http.Error(w, "Failed to merge page content", http.StatusInternalServerError)
				return
			}
		}

		reqCont.LogHistory = append(reqCont.LogHistory, "Response code: "+strconv.Itoa(interceptor.statusCode))

		if modifiedContent == nil {
			modifiedContent = body_node
		}

		bytesContent, err := parser.HTMLToByte(modifiedContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to convert HTML to byte: %v\n", err.Error())
			http.Error(w, "Failed to convert HTML to byte", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(bytesContent)))
		w.WriteHeader(interceptor.statusCode)
		if _, err := w.Write(bytesContent); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})
}
