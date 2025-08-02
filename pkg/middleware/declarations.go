package middleware

import (
	"bytes"
	"homeserver/pkg/globals"
	"homeserver/pkg/parser"
	"homeserver/pkg/schemas/user_schema"
	"homeserver/pkg/templates"
	"homeserver/pkg/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func MustParseFiles(filenames ...string) *template.Template {
	tmpl, err := template.ParseFiles(filenames...)
	if err != nil {
		log.Fatalf("Failed to parse template files: %v", err)
	}
	return tmpl
}

// Global variables for mocking
var (
	ParseFiles            = MustParseFiles
	IconsInjector         = utils.IconsInjector
	Must                  = utils.Must
	ParseTemplateToString = parser.ParseTemplateToString
	ErrorResponse         = utils.ErrorResponse
	Now                   = time.Now
	Since                 = time.Since
	IsProtectedPath       = globals.IsProtectedPath
)

// Context key for storing AdditionalRequestContext in request context
type contextKey string

const ReqContKey contextKey = "add_req_context"

// LogHistoryInterface defines logging methods
type LogHistoryInterface interface {
	Log(entry string)
	ErrLog(entry, at string)
	UnloadLogs(header string)
}

// LogHistory implements LogHistoryInterface
type LogHistory struct {
	Entries []string
}

func (lh *LogHistory) Log(entry string) {
	lh.Entries = append(lh.Entries, entry)
}

func (lh *LogHistory) ErrLog(entry, at string) {
	lh.Entries = append(lh.Entries, "[Err at "+at+"] "+entry)
}

func (lh *LogHistory) UnloadLogs(header string) {
	logs := strings.Join(lh.Entries, "\n") + "\n\n"
	log.Println("\n\n" + header + logs)
}

// AdditionalRequestContext holds request-specific data
type AdditionalRequestContext struct {
	LogHistory    LogHistoryInterface
	ToRestorePage bool
	RerenderBody  bool
	AuthToken     string
	Requester     *user_schema.User
	PageContent   *templates.LayoutPageData
}

// Default context initialization
var defaultContext = &AdditionalRequestContext{
	LogHistory:    &LogHistory{Entries: []string{}},
	ToRestorePage: false,
	RerenderBody:  false,
	AuthToken:     "",
	Requester:     nil,
	PageContent:   &templates.LayoutPageData{},
}

// Response interceptor for capturing response details
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

// WriterWithIcons injects icons into HTML responses
type WriterWithIcons struct {
	http.ResponseWriter
}

func (w *WriterWithIcons) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		return w.ResponseWriter.Write(b)
	}

	b = IconsInjector(b)

	return w.ResponseWriter.Write(b)
}
