package middleware

import (
	"bytes"
	"fmt"
	"homeserver/pkg/parser"
	"html/template"
	"net/http"
	"os"
	"strings"

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
	if ri.ResponseWriter.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		return ri.body.Write(b)
	}

	b = InterceptIcons(b)

	return ri.body.Write(b)
}

func Interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCont, ok := r.Context().Value(ReqContKey).(*AdditionalRequestContext)
		if !ok {
			fmt.Fprintf(os.Stderr, "Missing request context in interceptor middleware\n")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		interceptor := InterceptorConstructor(w)

		next.ServeHTTP(interceptor, r)

		reqCont.LogHistory.Log("[Interceptor] Response status code: " + fmt.Sprint(interceptor.statusCode))

		if interceptor.statusCode >= 400 {
			reqCont.LogHistory.ErrLog("[Interceptor] Response status code indicates an error: "+fmt.Sprint(interceptor.statusCode), "Interceptor:50")

			http.Error(w, http.StatusText(interceptor.statusCode), interceptor.statusCode)
			return
		}

		if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
			if _, err := w.Write(interceptor.body.Bytes()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write response: %v\n", err.Error())
			}
			return
		}

		if reqCont.ToRestorePage {
			reqCont.LogHistory.Log("[Interceptor] ToRestorePage is true, restoring page")
			tmpl := template.Must(template.ParseFiles("../static/index.html", "../static/layout.html"))

			reqCont.PageContent.Content = template.HTML(interceptor.body.String())

			if err := tmpl.ExecuteTemplate(w, "index.html", reqCont.PageContent); err != nil {
				reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "Interceptor:89")
				http.Error(w, "Failed to render page", http.StatusInternalServerError)
			}
			return
		}

		if reqCont.RerenderBody {
			reqCont.LogHistory.Log("[Interceptor] RerenderBody is true, rendering body")
			w.Header().Set("HX-Retarget", "body")
			w.Header().Set("HX-Reswap", "innerHTML")

			tmpl := template.Must(template.ParseFiles("../static/layout.html"))

			reqCont.PageContent.Content = template.HTML(interceptor.body.String())

			if err := tmpl.Execute(w, reqCont.PageContent); err != nil {
				reqCont.LogHistory.ErrLog("Failed to execute template: "+err.Error(), "Interceptor:60")
				http.Error(w, "Failed to render page", http.StatusInternalServerError)
			}
			return
		}

		body_node, err := parser.ParseString(interceptor.body.String())
		if err != nil {
			reqCont.LogHistory.ErrLog("Failed to parse response body: "+err.Error()+"\nBody: "+interceptor.body.String(), "Interceptor:65")
			http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
			return
		}

		if err := html.Render(w, body_node); err != nil {
			reqCont.LogHistory.ErrLog("Failed to render HTML response: "+err.Error(), "Interceptor:86")
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})
}

func InterceptIcons(b []byte) []byte {
	str := string(b)

	// fmt.Printf("Intercepting icons in response: %s\n", str)

	for i := 0; i <= len(str)-7; i++ {
		if str[i:i+6] != "<icon " {
			continue
		}

		cut_str := str[i:]

		// fmt.Printf("Found icon tag at index %d\n", i)

		end_index := strings.Index(cut_str, "</icon>")
		if end_index == -1 {
			fmt.Fprintf(os.Stderr, "Malformed icon tag in response\n")
			continue
		}

		end_index += len("</icon>")

		cut_str = cut_str[:end_index]
		// fmt.Printf("Cut string for icon tag: %s\n", cut_str)

		// fmt.Printf("End of icon tag found at index %d\n", end_index)

		name_start_idx := strings.Index(cut_str, "\"")
		if name_start_idx == -1 {
			fmt.Fprintf(os.Stderr, "Missing name attribute in icon tag\n")
			continue
		}

		name_end_idx := strings.Index(cut_str[name_start_idx+1:], "\"")
		if name_end_idx == -1 {
			fmt.Fprintf(os.Stderr, "Malformed icon tag in response\n")
			continue
		}

		name_end_idx += name_start_idx + 1

		iconName := cut_str[name_start_idx+1 : name_end_idx]
		// fmt.Printf("Found icon: %s\n", iconName)

		iconSVG := parser.LoadIcon(iconName)

		str = str[:i] + iconSVG + str[i+end_index:] // 7 is the length of "</icon>"
		i += len(iconSVG) - 1                       // Adjust index to account for the new length
	}

	return []byte(str)
}
